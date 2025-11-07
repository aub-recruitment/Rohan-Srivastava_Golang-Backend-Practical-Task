package handlers

import (
	"net/http"
	"strconv"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContentHandler struct {
	contentUseCase *usecases.ContentUseCase
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
	c.JSON(http.StatusOK, content)
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
	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
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
