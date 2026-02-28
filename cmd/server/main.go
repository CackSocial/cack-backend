package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/CackSocial/cack-backend/docs"
	"github.com/CackSocial/cack-backend/internal/handler"
	"github.com/CackSocial/cack-backend/internal/handler/ws"
	"github.com/CackSocial/cack-backend/internal/infrastructure/database"
	"github.com/CackSocial/cack-backend/internal/infrastructure/database/repository"
	"github.com/CackSocial/cack-backend/internal/infrastructure/storage"
	"github.com/CackSocial/cack-backend/internal/middleware"
	commentUC "github.com/CackSocial/cack-backend/internal/usecase/comment"
	followUC "github.com/CackSocial/cack-backend/internal/usecase/follow"
	likeUC "github.com/CackSocial/cack-backend/internal/usecase/like"
	messageUC "github.com/CackSocial/cack-backend/internal/usecase/message"
	postUC "github.com/CackSocial/cack-backend/internal/usecase/post"
	tagUC "github.com/CackSocial/cack-backend/internal/usecase/tag"
	timelineUC "github.com/CackSocial/cack-backend/internal/usecase/timeline"
	userUC "github.com/CackSocial/cack-backend/internal/usecase/user"
	bookmarkUC "github.com/CackSocial/cack-backend/internal/usecase/bookmark"
	"github.com/CackSocial/cack-backend/pkg/config"
)

// @title SocialConnect API
// @version 1.0
// @description Social networking platform API with posts, follows, messaging, and more.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}" to authenticate
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db := database.NewPostgresDB(cfg)

	// Create upload directory
	if err := os.MkdirAll(cfg.UploadPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	tagRepo := repository.NewTagRepository(db)
	followRepo := repository.NewFollowRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	// Initialize storage
	localStorage := storage.NewLocalStorage(cfg.UploadPath, cfg.BaseURL)

	// Initialize use cases
	userUseCase := userUC.NewUserUseCase(userRepo, followRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	postUseCase := postUC.NewPostUseCase(postRepo, tagRepo, likeRepo, commentRepo, userRepo, bookmarkRepo, localStorage)
	followUseCase := followUC.NewFollowUseCase(followRepo, userRepo)
	timelineUseCase := timelineUC.NewTimelineUseCase(followRepo, postRepo, likeRepo, commentRepo, bookmarkRepo)
	likeUseCase := likeUC.NewLikeUseCase(likeRepo, postRepo, userRepo)
	commentUseCase := commentUC.NewCommentUseCase(commentRepo, postRepo)
	messageUseCase := messageUC.NewMessageUseCase(messageRepo, userRepo, localStorage)
	tagUseCase := tagUC.NewTagUseCase(tagRepo, postRepo, likeRepo, commentRepo, bookmarkRepo)
	bookmarkUseCase := bookmarkUC.NewBookmarkUseCase(bookmarkRepo, postRepo, likeRepo, commentRepo, userRepo)

	// Initialize WebSocket hub and start it
	hub := ws.NewHub(messageUseCase)
	go hub.Run()

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)
	postHandler := handler.NewPostHandler(postUseCase)
	followHandler := handler.NewFollowHandler(followUseCase)
	timelineHandler := handler.NewTimelineHandler(timelineUseCase)
	likeHandler := handler.NewLikeHandler(likeUseCase)
	commentHandler := handler.NewCommentHandler(commentUseCase)
	messageHandler := handler.NewMessageHandler(messageUseCase)
	tagHandler := handler.NewTagHandler(tagUseCase)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkUseCase)
	wsHandler := ws.NewWSHandler(hub, cfg.JWTSecret)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	// Static file serving
	router.Static("/uploads", cfg.UploadPath)

	// API v1 route groups
	api := router.Group("/api/v1")

	public := api.Group("")
	optionalAuth := middleware.OptionalAuth(cfg.JWTSecret)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Register routes
	userHandler.RegisterRoutes(public, protected, optionalAuth)
	postHandler.RegisterRoutes(public, protected, optionalAuth)
	followHandler.RegisterRoutes(public, protected)
	timelineHandler.RegisterRoutes(protected)
	likeHandler.RegisterRoutes(public, protected)
	commentHandler.RegisterRoutes(public, protected)
	messageHandler.RegisterRoutes(protected)
	tagHandler.RegisterRoutes(public, optionalAuth)
	bookmarkHandler.RegisterRoutes(protected)
	wsHandler.RegisterRoutes(router)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
