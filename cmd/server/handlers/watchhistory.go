package handlers

import (
	"net/http"
	"strconv"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WatchHistoryHandler struct {
	watchHistoryUseCase *usecases.WatchHistoryUseCase
}

// @name NewWatchHistoryHandler - Creates new instance of watch history handler
// @param watchHistoryUseCase - watch history service instance
// @returns - new instance of watch history handler
func NewWatchHistoryHandler(watchHistoryUseCase *usecases.WatchHistoryUseCase) *WatchHistoryHandler {
	return &WatchHistoryHandler{watchHistoryUseCase: watchHistoryUseCase}
}

// @name CreateOrUpdateWatchHistory - Create/Update user's watch history
// @param c - gin context
// @returns - New watch history
// @dev - checks access level, use watched seconds = 0 for checking access
//
//	for a given content_id
func (h *WatchHistoryHandler) CreateOrUpdateWatchHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	var input usecases.WatchHistoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	watchHistory, err := h.watchHistoryUseCase.CreateOrUpdateWatchHistory(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, watchHistory)
}

// @name GetWatchHistory - Get's user's watch history
// @param c - gin context
// @returns - user's recent watch history in desc order of time
func (h *WatchHistoryHandler) GetWatchHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	histories, total, err := h.watchHistoryUseCase.GetWatchHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"watch_history": histories,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}

// @name GetContinueWatching - Get's user's unfinished content
// @param c - gin context
// @returns - user's recent unfinished content in desc order of time
func (h *WatchHistoryHandler) GetContinueWatching(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	histories, err := h.watchHistoryUseCase.GetContinueWatching(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"continue_watching": histories})
}

// @name UpdateProgress - Updates user's progress on a given content
// @param c - gin context
// @returns - newly updated watch history for single content
func (h *WatchHistoryHandler) UpdateProgress(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	historyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid history ID"})
		return
	}
	var input struct {
		WatchedSeconds *int `json:"watched_seconds" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	watchHistory, err := h.watchHistoryUseCase.UpdateProgress(c.Request.Context(), userID, historyID, *input.WatchedSeconds)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, watchHistory)
}
