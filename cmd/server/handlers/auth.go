package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

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
