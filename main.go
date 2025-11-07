package main

import (
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/controllers"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
}

func main() {
	router := gin.Default()

	// MIDDLEWARE
	router.Use(middleware.CustomResponseMiddleware())

	// USER ROUTES
	// CONTENT ROUTES
	// SUBSCRIPTION ROUTES

	// DEFAULT ROUTES
	router.GET("/v1/health", controllers.HealthCheck)
	router.NoRoute(middleware.NoRouteMiddleware())

	router.Run()
}
