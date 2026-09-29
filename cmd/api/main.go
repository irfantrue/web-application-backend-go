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
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := deliveryHttp.NewFiber(cfg)

	healthUsecase := usecase.NewHealthUsecase()
	healthHandler := handler.NewHealthHandler(healthUsecase)

	setupRoutes := &deliveryHttp.SetupRoutes{
		App:           app,
		HealthHandler: healthHandler,
	}

	setupRoutes.Build()

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
		log.Fatalf("Fiber shutdown failed.")
	}

	// Close other external connections sequentially.

	log.Println("Graceful shutdown complete. The application has stopped safely.")
}
