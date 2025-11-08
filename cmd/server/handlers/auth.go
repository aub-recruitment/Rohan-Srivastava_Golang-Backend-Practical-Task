package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

// NewAuthHandler - Sets up a new instance of the AuthHandler
// @param authUseCase - auth service containing business logic
// @returns - AuthHandler instance
func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// @name Register - Register a new user
// @param c - gin context
// @returns - access_token and refresh_token
func (h *AuthHandler) Register(c *gin.Context) {
	var input usecases.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authUseCase.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

// @name Login - Logs in an existing user
// @param c - gin context
// @returns - access_token and refresh_token
func (h *AuthHandler) Login(c *gin.Context) {
	var input usecases.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authUseCase.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @name Refresh - Given existing refresh_token generate new JWTs
// @param c - gin context
// @returns - access_token and refresh_token
func (h *AuthHandler) Refresh(c *gin.Context) {
	userID := c.GetString("userID")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	response, err := h.authUseCase.Refresh(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
