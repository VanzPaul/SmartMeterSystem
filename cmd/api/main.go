package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"SmartMeterSystem/internal"
	"SmartMeterSystem/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	logger, loggerErr := internal.NewLogger()
	if loggerErr != nil {
		panic(loggerErr)
	}
	defer logger.Sync()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Sugar().Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Sugar().Infof("Server forced to shutdown with error: %v", err)
	}

	logger.Sugar().Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	logger, loggerErr := internal.NewLogger()
	if loggerErr != nil {
		panic(loggerErr)
	}
	defer logger.Sync()

	// New server
	server := server.NewServer()

	// Construct the full address
	fullAddress := fmt.Sprintf("http://%s%s", internal.GetResolvedIP(), server.Addr)
	logger.Sugar().Infof("Server is running at %s\n", fullAddress)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	serverErr := server.ListenAndServe()
	if serverErr != nil && serverErr != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", serverErr))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
