package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/routes"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.uber.org/zap"
)

func main() {
	// Initialize thread-safe logger
	// config.InitLogger()
	// defer utils.Logger.Sync()

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		utils.Logger.Fatal("Failed to load environment variables",
			zap.Error(err))
	}

	// Create HTTP server with routes and middleware
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes.ServeMuxInit(), // Contains all global middleware
	}

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		utils.Logger.Info("🚀 Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Server failed to start",
				zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-stop
	utils.Logger.Info("🛑 Received shutdown signal")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful server shutdown
	utils.Logger.Info("⏳ Shutting down gracefully...")
	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.Error("Server shutdown error",
			zap.Error(err))
	}

	utils.Logger.Info("✅ Server stopped gracefully")
}
