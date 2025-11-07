package controllers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/utils"
	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	ctx.Set(utils.ResponseDataKey, gin.H{"message": "working as intended..."})
	ctx.Status(http.StatusOK)
}
