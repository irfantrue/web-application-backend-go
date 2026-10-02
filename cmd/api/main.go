// Package main is the entry point of the application.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	deliveryHttp "github.com/irfantrue/web-application-backend-go/internal/delivery/http"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/handler"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/worker"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/broker"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/cache"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/database"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/repository"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/security"
	cryptoService "github.com/irfantrue/web-application-backend-go/internal/pkg/crypto"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	// Database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Printf("Database error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB: %v", err)
	}
	defer func() {
		log.Println("Closing Postgres connection...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database: %v\n", err)
		}
	}()

	// Redis (optional — failure is not critical)
	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Printf("Redis error: %v", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v\n", err)
		}
	}()

	// RabbitMQ
	broker, err := broker.NewRabbitMQ(cfg)
	if err != nil {
		log.Printf("Failed to initialize RabbitMQ: %v", err)
	}

	defer func() {
		if err := broker.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ: %v", err)
		}
	}()

	// Worker — non-blocking,
	go worker.StartNotificationWorker(ctx, broker)

	// Dependency wiring
	uow := database.NewGormUnitOfWork(db)

	cryptoService := cryptoService.NewCryptographyService()
	tokenService := security.NewJwtService(cfg)
	cacheService := cache.NewRedisCache(redisClient)

	apiKeyRepository := repository.NewAPIKeyRepository(db)
	userRepository := repository.NewUserRepository(db)
	refreshTokenSessionRepository := repository.NewRefreshTokenSessionRepository(db)

	apiKeyUsecase := usecase.NewAPIKeyUsecase(uow, apiKeyRepository, cacheService, cryptoService)
	authUsecase := usecase.NewAuthUsecase(uow, userRepository, refreshTokenSessionRepository, cryptoService, tokenService, broker)

	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)

	app := deliveryHttp.NewFiber(cfg)
	setupRoutes := &deliveryHttp.Config{
		App:           app,
		TokenService:  tokenService,
		APIKeyService: apiKeyUsecase,
		Config:        cfg,
		APIKeyHandler: apiKeyHandler,
		AuthHandler:   authHandler,
	}
	deliveryHttp.NewSetupRoutes(setupRoutes)

	// HTTP server — non-blocking
	go func() {
		log.Printf("Server is running on port %s", cfg.App.Port)
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for ctx to be cancelled (SIGINT/SIGTERM)
	<-ctx.Done()
	log.Println("Received termination signal. Initiating graceful shutdown...")

	// Graceful HTTP shutdown (allow time to drain active requests)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Fiber shutdown failed.")
	}

	// Close other external connections sequentially.

	log.Println("Graceful shutdown complete. The application has stopped safely.")
}
