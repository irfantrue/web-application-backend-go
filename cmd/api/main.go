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
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/cache"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/database"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/repository"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	defer func() {
		log.Println("Closing Postgres connection...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database: %v\n", err)
		}
	}()

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Printf("Redis error: %v", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v\n", err)
		}
	}()

	uow := database.NewGormUnitOfWork(db)
	productRepository := repository.NewProductRepository(db)

	_ = cache.NewRedisCache(redisClient)

	productUsecase := usecase.NewProductUsecase(uow, productRepository)
	productHandler := handler.NewProductHandler(productUsecase)

	app := deliveryHttp.NewFiber(cfg)

	setupRoutes := &deliveryHttp.Config{
		App:            app,
		Config:         cfg,
		ProductHandler: productHandler,
	}

	deliveryHttp.NewSetupRoutes(setupRoutes)

	go func() {
		log.Printf("Server is running on port %s", cfg.App.Port)
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	log.Println("Received termination signal. Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Fiber shutdown failed.")
	}

	// Close other external connections sequentially.

	log.Println("Graceful shutdown complete. The application has stopped safely.")
}
