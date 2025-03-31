package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HumbleBee14/distributed-log-aggregator/internal/api"
	"github.com/HumbleBee14/distributed-log-aggregator/internal/config"
	"github.com/HumbleBee14/distributed-log-aggregator/internal/storage"
	"github.com/HumbleBee14/distributed-log-aggregator/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	logger.InitFromEnv()

	// Load configuration
	cfg := config.New()

	// logger.Info("Configuration loaded successfully")

	// Initialize Redis storage
	redisStorage, err := initializeStorage(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize storage: %v", err)
	}
	defer redisStorage.Close()

	// Set up the router
	router := api.NewRouter(redisStorage)
	logger.Info("API router initialized")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router.Setup(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("Starting server on port %s", cfg.Server.Port)
		serverErr <- server.ListenAndServe()
	}()

	// Set up graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt or server error
	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	case <-interrupt:
		logger.Info("Received interrupt signal, shutting down gracefully...")

		// Create a deadline for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server shutdown failed: %v", err)
		}
		logger.Info("Server gracefully stopped")
	}
}

func initializeStorage(cfg *config.Config) (*storage.RedisStorage, error) {
	const maxRetries = 4
	const retryDelay = 2 * time.Second

	var err error
	var redisStore *storage.RedisStorage

	for attempt := 1; attempt <= maxRetries; attempt++ {
		redisStore, err = storage.NewRedisStorage(&cfg.Redis)
		if err == nil {
			logger.Info("Storage connection established")
			return redisStore, nil
		}

		if attempt < maxRetries {
			logger.Warn("Storage connection failed (attempt %d/%d): %v",
				attempt, maxRetries, err)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to storage after %d attempts: %w",
		maxRetries, err)
}
