package notify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type WebhookConfig struct {
	URL     string
	Secret  string
	Headers map[string]string
	Timeout time.Duration
}

type WebhookSender struct {
	config *WebhookConfig
	client *http.Client
}

func NewWebhookSender(cfg *WebhookConfig) *WebhookSender {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &WebhookSender{
		config: cfg,
		client: &http.Client{Timeout: timeout},
	}
}

type WebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type WebhookAttachment struct {
	Filename    string
	ContentType string
	Reader      io.Reader
}

//nolint:gocyclo // Webhook sending has inherent complexity with multipart, retries, HMAC
func (w *WebhookSender) Send(ctx context.Context, payload *WebhookPayload, attachments []WebhookAttachment) error {
	var body bytes.Buffer
	var contentType string

	if len(attachments) > 0 {
		writer := multipart.NewWriter(&body)

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		payloadField, err := writer.CreateFormField("payload")
		if err != nil {
			return fmt.Errorf("failed to create payload field: %w", err)
		}
		if _, err := payloadField.Write(payloadJSON); err != nil {
			return fmt.Errorf("failed to write payload: %w", err)
		}

		for _, att := range attachments {
			part, err := writer.CreateFormFile("files", att.Filename)
			if err != nil {
				return fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := io.Copy(part, att.Reader); err != nil {
				return fmt.Errorf("failed to copy attachment: %w", err)
			}
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close multipart writer: %w", err)
		}

		contentType = writer.FormDataContentType()
	} else {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return fmt.Errorf("failed to encode payload: %w", err)
		}
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.config.URL, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Pagemail-Webhook/1.0")

	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	if w.config.Secret != "" {
		signature := w.computeSignature(body.Bytes())
		req.Header.Set("X-Pagemail-Signature", signature)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (w *WebhookSender) computeSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(w.config.Secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func (w *WebhookSender) Test(ctx context.Context) error {
	payload := &WebhookPayload{
		Event:     "test",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: map[string]interface{}{
			"message": "This is a test webhook from Pagemail",
		},
	}

	return w.Send(ctx, payload, nil)
}
