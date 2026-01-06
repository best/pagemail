package handlers

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"pagemail/internal/audit"
	"pagemail/internal/models"
	"pagemail/internal/pkg/avatar"
	"pagemail/internal/pkg/errors"
)

const defaultContentType = "application/octet-stream"

func handleFormFileError(c *gin.Context, err error) {
	var maxBytesErr *http.MaxBytesError
	if stderrors.As(err, &maxBytesErr) {
		errors.NewProblemDetail(http.StatusRequestEntityTooLarge, "Payload Too Large",
			fmt.Sprintf("File too large, max %d bytes", avatar.MaxUploadBytes)).Respond(c)
		return
	}
	if stderrors.Is(err, http.ErrMissingFile) {
		errors.BadRequest("Missing avatar file").Respond(c)
		return
	}
	errors.BadRequest("Failed to read upload").Respond(c)
}

func handleAvatarProcessError(c *gin.Context, err error) {
	switch {
	case stderrors.Is(err, avatar.ErrFileTooLarge):
		errors.NewProblemDetail(http.StatusRequestEntityTooLarge, "Payload Too Large", err.Error()).Respond(c)
	case stderrors.Is(err, avatar.ErrUnsupportedType):
		errors.NewProblemDetail(http.StatusUnsupportedMediaType, "Unsupported Media Type", err.Error()).Respond(c)
	default:
		errors.BadRequest(err.Error()).Respond(c)
	}
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		errors.BadRequest("Invalid user ID").Respond(c)
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", uid).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	maxBodyBytes := int64(avatar.MaxUploadBytes) + 1<<20
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		handleFormFileError(c, err)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		errors.BadRequest("Failed to read file").Respond(c)
		return
	}
	defer f.Close()

	result, err := avatar.ValidateAndProcess(f, avatar.MaxUploadBytes)
	if err != nil {
		handleAvatarProcessError(c, err)
		return
	}

	key := avatar.BuildStorageKey(uid, result.Extension)

	_, err = h.storage.Upload(c.Request.Context(), key, bytes.NewReader(result.Data), result.ContentType)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to upload avatar")
		errors.InternalError("Failed to store avatar").Respond(c)
		return
	}

	oldKey := user.AvatarKey
	user.AvatarKey = &key
	if err := h.db.Save(&user).Error; err != nil {
		_ = h.storage.Delete(c.Request.Context(), key)
		errors.InternalError("Failed to update user").Respond(c)
		return
	}

	if oldKey != nil && *oldKey != "" && *oldKey != key {
		if err := h.storage.Delete(c.Request.Context(), *oldKey); err != nil {
			log.Warn().Err(err).Str("key", *oldKey).Msg("Failed to delete old avatar")
		}
	}

	h.logAudit(c, audit.ActionUserUpdate, "user", &user.ID, audit.ResourceDetails{
		Email: user.Email,
	})

	c.JSON(http.StatusOK, buildUserResponse(&user))
}

func (h *Handler) GetAvatar(c *gin.Context) {
	targetID := c.Param("id")
	if targetID == "me" {
		targetID = c.GetString("user_id")
	}

	uid, err := uuid.Parse(targetID)
	if err != nil {
		errors.BadRequest("Invalid user ID").Respond(c)
		return
	}

	currentUserID := c.GetString("user_id")
	currentRole := c.GetString("user_role")

	if currentUserID != targetID && currentRole != models.RoleAdmin {
		errors.Forbidden("Access denied").Respond(c)
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", uid).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	if user.AvatarKey == nil || *user.AvatarKey == "" {
		errors.NotFound("No avatar").Respond(c)
		return
	}

	reader, info, err := h.storage.Download(c.Request.Context(), *user.AvatarKey)
	if err != nil {
		log.Error().Err(err).Str("key", *user.AvatarKey).Msg("Failed to download avatar")
		errors.NotFound("Avatar not found").Respond(c)
		return
	}
	defer reader.Close()

	contentType := info.ContentType
	if contentType == "" || contentType == defaultContentType {
		if ct := mime.TypeByExtension(path.Ext(*user.AvatarKey)); ct != "" {
			contentType = ct
		} else {
			contentType = defaultContentType
		}
	}

	c.Header("Content-Type", contentType)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "private, max-age=3600")
	if info.Size > 0 {
		c.Header("Content-Length", fmt.Sprintf("%d", info.Size))
	}
	c.Status(http.StatusOK)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		log.Error().Err(err).Msg("Failed to stream avatar")
	}
}

func (h *Handler) DeleteAvatar(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		errors.BadRequest("Invalid user ID").Respond(c)
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", uid).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	if user.AvatarKey == nil || *user.AvatarKey == "" {
		c.JSON(http.StatusOK, buildUserResponse(&user))
		return
	}

	oldKey := *user.AvatarKey
	user.AvatarKey = nil
	if err := h.db.Save(&user).Error; err != nil {
		errors.InternalError("Failed to update user").Respond(c)
		return
	}

	if err := h.storage.Delete(c.Request.Context(), oldKey); err != nil {
		log.Warn().Err(err).Str("key", oldKey).Msg("Failed to delete avatar file")
	}

	h.logAudit(c, audit.ActionUserUpdate, "user", &user.ID, audit.ResourceDetails{
		Email: user.Email,
	})

	c.JSON(http.StatusOK, buildUserResponse(&user))
}

func buildUserResponse(user *models.User) gin.H {
	resp := gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	if user.AvatarKey != nil && *user.AvatarKey != "" {
		resp["avatar_url"] = fmt.Sprintf("/users/%s/avatar?t=%d", user.ID, user.UpdatedAt.Unix())
	}
	return resp
}
