package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContentHandler struct {
	contentUseCase *usecases.ContentUseCase
}

// ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
// 	Title           string      `gorm:"not null;index" json:"title"`
// 	Description     string      `json:"description"`
// 	AccessLevel     AccessLevel `gorm:"type:varchar(20);not null;index" json:"access_level"`
// 	DurationSeconds int         `gorm:"not null" json:"duration_seconds"`
// 	ThumbnailURL    string      `json:"thumbnail_url"`
// 	TrailerURL      string      `json:"trailer_url"`
// 	VideoURL        string      `json:"video_url"`
// 	Published       bool        `gorm:"default:false;index" json:"published"`
// 	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
// 	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

type OpenContentOutput struct {
	ID           uuid.UUID           `json:"id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	AccessLevel  *domain.AccessLevel `json:"access_level"`
	Duration     int                 `json:"duration"`
	ThumbnailURL string              `json:"thumbnail_url"`
	TrailerURL   string              `json:"trailer_url"`
	CreatedAt    time.Time           `json:"published_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func NewContentHandler(contentUseCase *usecases.ContentUseCase) *ContentHandler {
	return &ContentHandler{contentUseCase: contentUseCase}
}

func (h *ContentHandler) CreateContent(c *gin.Context) {
	var input usecases.CreateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := h.contentUseCase.CreateContent(c.Request.Context(), input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, content)
}

func (h *ContentHandler) GetContent(c *gin.Context) {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content ID"})
		return
	}
	var userID *uuid.UUID
	if userIDStr := c.GetString("userID"); userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err == nil {
			userID = &id
		}
	}
	content, err := h.contentUseCase.GetContent(c.Request.Context(), contentID, userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, OpenContentOutput{
		ID:           content.ID,
		Title:        content.Title,
		Description:  content.Description,
		AccessLevel:  &content.AccessLevel,
		Duration:     content.DurationSeconds,
		ThumbnailURL: content.ThumbnailURL,
		TrailerURL:   content.TrailerURL,
		CreatedAt:    content.CreatedAt,
		UpdatedAt:    content.UpdatedAt,
	})
}

func (h *ContentHandler) ListContent(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	publishedOnly := c.DefaultQuery("published", "true") == "true"
	contents, total, err := h.contentUseCase.ListContent(c.Request.Context(), publishedOnly, limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	outputs := make([]OpenContentOutput, len(contents))
	for i, c := range contents {
		outputs[i] = OpenContentOutput{
			ID:           c.ID,
			Title:        c.Title,
			Description:  c.Description,
			AccessLevel:  &c.AccessLevel,
			Duration:     c.DurationSeconds,
			ThumbnailURL: c.ThumbnailURL,
			TrailerURL:   c.TrailerURL,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": outputs,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ContentHandler) UpdateContent(c *gin.Context) {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content ID"})
		return
	}
	var input usecases.CreateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := h.contentUseCase.UpdateContent(c.Request.Context(), contentID, input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (h *ContentHandler) DeleteContent(c *gin.Context) {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content ID"})
		return
	}
	if err := h.contentUseCase.DeleteContent(c.Request.Context(), contentID); err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "content deleted successfully"})
}
