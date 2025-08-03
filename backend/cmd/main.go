package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/xarcher/backend/config"
	"github.com/xarcher/backend/internal/delivery/handler"
	"github.com/xarcher/backend/internal/delivery/handler/middleware"
	"github.com/xarcher/backend/internal/infrastructure/database"
	"github.com/xarcher/backend/internal/infrastructure/jwt"
	"github.com/xarcher/backend/internal/repository"
	"github.com/xarcher/backend/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Database connection
	dbConfig := database.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(db)

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create temp directory for uploads if not exists
	if err := os.MkdirAll(cfg.Upload.TempDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Services
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey)

	// Repositories
	userRepository := repository.NewUserRepository(db)
	uploadRepository := repository.NewUploadRepository(db)

	// Use cases
	authUsecase := usecase.NewAuthUsecase(userRepository, jwtService, 10*time.Second)
	uploadUsecase := usecase.NewUploadUsecase(uploadRepository, cfg.Upload, 10*time.Second)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	uploadHandler := handler.NewUploadHandler(uploadUsecase)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	// Routes
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/revoke", authHandler.RevokeToken).Methods("POST")

	// Upload routes
	r.HandleFunc("/upload-form", uploadHandler.ServeUploadForm).Methods("GET")
	r.HandleFunc("/upload", authMiddleware.Authenticate(uploadHandler.UploadFile)).Methods("POST")

	// CORS setup for development
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify exact origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	handler := c.Handler(r)

	// Server configuration
	srv := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", cfg.GetServerAddress())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
