package queue

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"pagemail/internal/capture"
	"pagemail/internal/config"
	"pagemail/internal/models"
	"pagemail/internal/storage"
)

type Dispatcher struct {
	cfg      *config.Config
	db       *gorm.DB
	storage  storage.Storage
	jobChan  chan models.Job
	workers  []*Worker
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	workerID string
}

func NewDispatcher(cfg *config.Config, db *gorm.DB, store storage.Storage) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Dispatcher{
		cfg:      cfg,
		db:       db,
		storage:  store,
		jobChan:  make(chan models.Job, 100),
		ctx:      ctx,
		cancel:   cancel,
		workerID: uuid.New().String()[:8],
	}
}

func (d *Dispatcher) Start() {
	log.Info().Int("workers", d.cfg.Capture.Workers).Msg("Starting job dispatcher")

	for i := 0; i < d.cfg.Capture.Workers; i++ {
		worker := NewWorker(i, d.cfg, d.db, d.storage, d.jobChan)
		d.workers = append(d.workers, worker)
		d.wg.Add(1)
		go func(w *Worker) {
			defer d.wg.Done()
			w.Start(d.ctx)
		}(worker)
	}

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.poll()
	}()

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.recoverStuckJobs()
	}()
}

func (d *Dispatcher) Stop() {
	log.Info().Msg("Stopping job dispatcher")
	d.cancel()
	close(d.jobChan)

	for _, w := range d.workers {
		w.Close()
	}

	d.wg.Wait()
	log.Info().Msg("Job dispatcher stopped")
}

func (d *Dispatcher) poll() {
	ticker := time.NewTicker(time.Duration(d.cfg.Queue.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.fetchAndDispatch()
		}
	}
}

func (d *Dispatcher) fetchAndDispatch() {
	var jobs []models.Job
	now := time.Now()
	leaseUntil := now.Add(time.Duration(d.cfg.Queue.LeaseDuration) * time.Second)

	result := d.db.
		Where("status = ? AND run_at <= ?", models.JobStatusPending, now).
		Order("priority DESC, run_at ASC").
		Limit(10).
		Find(&jobs)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to fetch jobs")
		return
	}

	//nolint:gocritic // rangeValCopy: job is sent to channel which requires value type
	for _, job := range jobs {
		lockedBy := d.workerID

		result := d.db.Model(&job).
			Where("status = ?", models.JobStatusPending).
			Updates(map[string]interface{}{
				"status":      models.JobStatusRunning,
				"locked_by":   lockedBy,
				"locked_at":   now,
				"lease_until": leaseUntil,
			})

		if result.RowsAffected > 0 {
			select {
			case d.jobChan <- job:
				log.Debug().Str("job_id", job.ID.String()).Str("type", job.Type).Msg("Job dispatched")
			default:
				d.db.Model(&job).Updates(map[string]interface{}{
					"status":      models.JobStatusPending,
					"locked_by":   nil,
					"locked_at":   nil,
					"lease_until": nil,
				})
				log.Warn().Str("job_id", job.ID.String()).Msg("Job channel full, returning job to queue")
			}
		}
	}
}

func (d *Dispatcher) recoverStuckJobs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			result := d.db.Model(&models.Job{}).
				Where("status = ? AND lease_until < ?", models.JobStatusRunning, now).
				Updates(map[string]interface{}{
					"status":      models.JobStatusPending,
					"locked_by":   nil,
					"locked_at":   nil,
					"lease_until": nil,
				})

			if result.RowsAffected > 0 {
				log.Info().Int64("count", result.RowsAffected).Msg("Recovered stuck jobs")
			}
		}
	}
}

type Worker struct {
	id      int
	cfg     *config.Config
	db      *gorm.DB
	storage storage.Storage
	jobChan <-chan models.Job
	browser *capture.Browser
	mu      sync.Mutex
}

func NewWorker(id int, cfg *config.Config, db *gorm.DB, store storage.Storage, jobChan <-chan models.Job) *Worker {
	return &Worker{
		id:      id,
		cfg:     cfg,
		db:      db,
		storage: store,
		jobChan: jobChan,
	}
}

func (w *Worker) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.browser != nil {
		w.browser.Close()
		w.browser = nil
	}
}

func (w *Worker) getBrowser() (*capture.Browser, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.browser != nil {
		return w.browser, nil
	}

	browser, err := capture.NewBrowser(&capture.BrowserConfig{
		Headless:       true,
		ViewportWidth:  w.cfg.Capture.ViewportWidth,
		ViewportHeight: w.cfg.Capture.ViewportHeight,
		Timeout:        time.Duration(w.cfg.Capture.WaitTimeout) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	w.browser = browser
	return w.browser, nil
}

func (w *Worker) Start(ctx context.Context) {
	log.Info().Int("worker_id", w.id).Msg("Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Int("worker_id", w.id).Msg("Worker stopped")
			return
		case job, ok := <-w.jobChan:
			if !ok {
				return
			}
			w.process(ctx, job)
		}
	}
}

//nolint:gocritic // hugeParam: job comes from channel which uses value type for simplicity
func (w *Worker) process(ctx context.Context, job models.Job) {
	log.Info().
		Int("worker_id", w.id).
		Str("job_id", job.ID.String()).
		Str("type", job.Type).
		Msg("Processing job")

	var err error
	switch job.Type {
	case models.JobTypeCapture:
		err = w.processCapture(ctx, job)
	case models.JobTypeDeliver:
		err = w.processDelivery(ctx, job)
	default:
		log.Error().Str("type", job.Type).Msg("Unknown job type")
		err = nil
	}

	if err != nil {
		w.handleFailure(&job, err)
	} else {
		w.handleSuccess(&job)
	}
}

type CapturePayload struct {
	TaskID  string   `json:"task_id"`
	URL     string   `json:"url"`
	Cookies string   `json:"cookies"`
	Formats []string `json:"formats"`
}

//nolint:gocyclo,gocritic // Capture processing has inherent complexity; job from channel uses value type
func (w *Worker) processCapture(ctx context.Context, job models.Job) error {
	log.Info().Str("job_id", job.ID.String()).Msg("Processing capture job")

	var payload CapturePayload
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	taskID, err := uuid.Parse(payload.TaskID)
	if err != nil {
		return fmt.Errorf("invalid task_id: %w", err)
	}

	var task models.CaptureTask
	if err := w.db.First(&task, "id = ?", taskID).Error; err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	w.db.Model(&task).Update("status", models.TaskStatusRunning)

	browser, err := w.getBrowser()
	if err != nil {
		w.updateTaskFailed(&task, fmt.Sprintf("browser init failed: %v", err))
		return fmt.Errorf("failed to get browser: %w", err)
	}

	cookies := parseCookies(payload.Cookies, payload.URL)

	opts := &capture.CaptureOptions{
		URL:            payload.URL,
		Cookies:        cookies,
		ViewportWidth:  task.ViewportWidth,
		ViewportHeight: task.ViewportHeight,
		UserAgent:      task.UserAgent,
		Timeout:        time.Duration(task.WaitTimeoutMs) * time.Millisecond,
	}

	log.Info().
		Str("task_id", taskID.String()).
		Str("url", payload.URL).
		Strs("formats", payload.Formats).
		Msg("Starting browser capture")

	result, err := browser.Capture(ctx, opts)
	if err != nil {
		w.updateTaskFailed(&task, fmt.Sprintf("capture failed: %v", err))
		return fmt.Errorf("capture failed: %w", err)
	}

	log.Info().
		Str("task_id", taskID.String()).
		Int("html_size", len(result.HTML)).
		Int("pdf_size", len(result.PDF)).
		Int("screenshot_size", len(result.Screenshot)).
		Msg("Capture completed, saving outputs")

	formatSet := make(map[string]bool)
	for _, f := range payload.Formats {
		formatSet[strings.ToLower(f)] = true
	}

	var outputs []models.CaptureOutput

	if formatSet["pdf"] && len(result.PDF) > 0 {
		output, err := w.saveOutput(ctx, taskID, "pdf", result.PDF, "application/pdf")
		if err != nil {
			log.Error().Err(err).Msg("Failed to save PDF")
		} else {
			outputs = append(outputs, *output)
		}
	}

	if formatSet["html"] && len(result.HTML) > 0 {
		output, err := w.saveOutput(ctx, taskID, "html", result.HTML, "text/html")
		if err != nil {
			log.Error().Err(err).Msg("Failed to save HTML")
		} else {
			outputs = append(outputs, *output)
		}
	}

	if formatSet["screenshot"] && len(result.Screenshot) > 0 {
		output, err := w.saveOutput(ctx, taskID, "screenshot", result.Screenshot, "image/png")
		if err != nil {
			log.Error().Err(err).Msg("Failed to save screenshot")
		} else {
			outputs = append(outputs, *output)
		}
	}

	if len(outputs) == 0 {
		w.updateTaskFailed(&task, "no outputs generated")
		return fmt.Errorf("no outputs generated")
	}

	if err := w.db.Create(&outputs).Error; err != nil {
		w.updateTaskFailed(&task, fmt.Sprintf("failed to save outputs: %v", err))
		return fmt.Errorf("failed to save outputs: %w", err)
	}

	now := time.Now()
	w.db.Model(&task).Updates(map[string]interface{}{
		"status":       models.TaskStatusCompleted,
		"completed_at": now,
	})

	log.Info().
		Str("task_id", taskID.String()).
		Int("output_count", len(outputs)).
		Msg("Capture task completed successfully")

	return nil
}

func (w *Worker) saveOutput(ctx context.Context, taskID uuid.UUID, format string, data []byte, contentType string) (*models.CaptureOutput, error) {
	ext := format
	if format == "screenshot" {
		ext = "png"
	}

	now := time.Now().UTC()
	// Path: captures/2025/12/01/20251201123022123456_uuid_format.ext (UTC)
	timestamp := now.Format("20060102150405") + fmt.Sprintf("%06d", now.Nanosecond()/1000)
	objectKey := fmt.Sprintf("captures/%s/%s_%s_%s.%s",
		now.Format("2006/01/02"),
		timestamp,
		taskID.String(),
		format,
		ext)

	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	info, err := w.storage.Upload(ctx, objectKey, bytes.NewReader(data), contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload %s: %w", format, err)
	}

	output := &models.CaptureOutput{
		TaskID:         taskID,
		Format:         format,
		StorageBackend: w.cfg.Storage.Backend,
		ObjectKey:      objectKey,
		ContentType:    contentType,
		SizeBytes:      info.Size,
		SHA256:         hashStr,
	}

	return output, nil
}

func (w *Worker) updateTaskFailed(task *models.CaptureTask, errMsg string) {
	w.db.Model(task).Updates(map[string]interface{}{
		"status":        models.TaskStatusFailed,
		"error_message": errMsg,
	})
}

func parseCookies(cookieStr, targetURL string) []*proto.NetworkCookieParam {
	if cookieStr == "" {
		return nil
	}

	pairs := strings.Split(cookieStr, ";")
	cookies := make([]*proto.NetworkCookieParam, 0, len(pairs))

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if name == "" {
			continue
		}

		cookies = append(cookies, &proto.NetworkCookieParam{
			Name:  name,
			Value: value,
			URL:   targetURL,
		})
	}

	return cookies
}

//nolint:gocritic,unparam // hugeParam: job from channel uses value type; ctx reserved for future use
func (w *Worker) processDelivery(ctx context.Context, job models.Job) error {
	log.Info().Str("job_id", job.ID.String()).Msg("Processing delivery job")
	return nil
}

func (w *Worker) handleSuccess(job *models.Job) {
	w.db.Model(job).Updates(map[string]interface{}{
		"status":      models.JobStatusSuccess,
		"locked_by":   nil,
		"locked_at":   nil,
		"lease_until": nil,
	})
	log.Info().Str("job_id", job.ID.String()).Msg("Job completed successfully")
}

func (w *Worker) handleFailure(job *models.Job, err error) {
	job.Attempts++
	job.LastError = err.Error()

	if job.Attempts >= job.MaxAttempts {
		w.db.Model(job).Updates(map[string]interface{}{
			"status":      models.JobStatusFailed,
			"attempts":    job.Attempts,
			"last_error":  job.LastError,
			"locked_by":   nil,
			"locked_at":   nil,
			"lease_until": nil,
		})
		log.Error().Str("job_id", job.ID.String()).Err(err).Msg("Job failed permanently")
	} else {
		retryDelay := time.Duration(10*(1<<job.Attempts)) * time.Second
		runAt := time.Now().Add(retryDelay)
		w.db.Model(job).Updates(map[string]interface{}{
			"status":      models.JobStatusPending,
			"attempts":    job.Attempts,
			"last_error":  job.LastError,
			"run_at":      runAt,
			"locked_by":   nil,
			"locked_at":   nil,
			"lease_until": nil,
		})
		log.Warn().Str("job_id", job.ID.String()).Err(err).Int("attempt", job.Attempts).Msg("Job failed, will retry")
	}
}

func EnqueueJob(db *gorm.DB, jobType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	job := models.Job{
		Type:        jobType,
		Payload:     string(payloadBytes),
		Status:      models.JobStatusPending,
		RunAt:       time.Now(),
		MaxAttempts: 3,
	}

	return db.Create(&job).Error
}
