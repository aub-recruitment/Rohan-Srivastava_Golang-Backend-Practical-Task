package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	user, err := h.userUseCase.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	var input usecases.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userUseCase.UpdateProfile(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetSubscriptionHistory(c *gin.Context) {
	userID := c.GetString("userID")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	subscriptions, err := h.userUseCase.GetSubscriptionHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}
