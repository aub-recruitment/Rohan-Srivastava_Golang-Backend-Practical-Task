package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/cmd/server/handlers"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/cmd/server/middleware"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories/postgres"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := infrastructure.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := infrastructure.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	cache := infrastructure.NewCache(cfg)
	if err := cache.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	jwtService := infrastructure.NewJWTService(cfg.JWTSecret, cfg.JWTSauce, cfg.JWTExpiration)

	userRepo := postgres.NewUserRepository(db)
	contentRepo := postgres.NewContentRepository(db)
	planRepo := postgres.NewPlanRepository(db)
	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	watchHistoryRepo := postgres.NewWatchHistoryRepository(db)

	authUseCase := usecases.NewAuthUseCase(userRepo, jwtService, cache, cfg.JWTExpiration)
	userUseCase := usecases.NewUserUseCase(userRepo, subscriptionRepo)
	contentUseCase := usecases.NewContentUseCase(contentRepo, subscriptionRepo, userRepo)
	planUseCase := usecases.NewPlanUseCase(planRepo)
	subscriptionUseCase := usecases.NewSubscriptionUseCase(subscriptionRepo, planRepo, userRepo)
	watchHistoryUseCase := usecases.NewWatchHistoryUseCase(watchHistoryRepo, contentRepo, subscriptionRepo)

	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)
	contentHandler := handlers.NewContentHandler(contentUseCase)
	planHandler := handlers.NewPlanHandler(planUseCase)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionUseCase)
	watchHistoryHandler := handlers.NewWatchHistoryHandler(watchHistoryUseCase)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	setupRoutes(router, cfg, jwtService, cache, authHandler, userHandler, contentHandler, planHandler, subscriptionHandler, watchHistoryHandler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func setupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	jwtService *infrastructure.JWTService,
	cache *infrastructure.Cache,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	contentHandler *handlers.ContentHandler,
	planHandler *handlers.PlanHandler,
	subscriptionHandler *handlers.SubscriptionHandler,
	watchHistoryHandler *handlers.WatchHistoryHandler,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.NoRoute(middleware.NoRouteMiddleware())

	v1 := router.Group("/api/v1")
	public := v1.Group("")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		content := public.Group("/content")
		{
			content.GET("", contentHandler.ListContent)
			content.GET("/:id", contentHandler.GetContent)
		}
		plans := public.Group("/plans")
		{
			plans.GET("", planHandler.ListPlans)
			plans.GET("/:id", planHandler.GetPlan)
		}
	}

	authMiddleware := middleware.AuthMiddleware(jwtService, cache)
	rateLimitMiddleware := middleware.RateLimitMiddleware(cache, int64(cfg.RequestLimit), time.Minute)

	protected := v1.Group("")
	protected.Use(authMiddleware)
	protected.Use(rateLimitMiddleware)
	{
		auth := protected.Group("/auth")
		{
			auth.GET("/refresh", authHandler.Refresh)
		}
		users := protected.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/subscription-history", userHandler.GetSubscriptionHistory)
		}
		subscriptions := protected.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("/active", subscriptionHandler.GetActiveSubscription)
			subscriptions.GET("/history", subscriptionHandler.GetSubscriptionHistory)
			subscriptions.POST("/:id/cancel", subscriptionHandler.CancelSubscription)
			subscriptions.POST("/:id/renew", subscriptionHandler.RenewSubscription)
		}
		watchHistory := protected.Group("/watch-history")
		{
			watchHistory.POST("", watchHistoryHandler.CreateOrUpdateWatchHistory)
			watchHistory.GET("", watchHistoryHandler.GetWatchHistory)
			watchHistory.GET("/continue-watching", watchHistoryHandler.GetContinueWatching)
			watchHistory.PUT("/:id", watchHistoryHandler.UpdateProgress)
		}
		adminMiddleware := middleware.AdminMiddleware()
		admin := protected.Group("/admin")
		admin.Use(adminMiddleware)
		{
			adminContent := admin.Group("/content")
			{
				adminContent.POST("", contentHandler.CreateContent)
				adminContent.PUT("/:id", contentHandler.UpdateContent)
				adminContent.DELETE("/:id", contentHandler.DeleteContent)
			}
			adminPlans := admin.Group("/plans")
			{
				adminPlans.POST("", planHandler.CreatePlan)
				adminPlans.PUT("/:id", planHandler.UpdatePlan)
				adminPlans.DELETE("/:id", planHandler.DeletePlan)
			}
		}
	}
}
