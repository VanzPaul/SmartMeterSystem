package routes

import (
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/middleware"
	"github.com/vanspaul/SmartMeterSystem/services"
)

func ServeMuxInit() http.Handler {
	// Create base router
	baseMux := http.NewServeMux()

	// Public routes (no auth required)
	baseMux.HandleFunc("/", services.Home)

	baseMux.HandleFunc("/register", services.Register)
	baseMux.HandleFunc("POST /register/submit", services.SubmitRegister)

	baseMux.HandleFunc("/login", services.Login)
	baseMux.HandleFunc("POST /login/submit", services.SubmitLogin)

	// Protected routes subsystem
	socketMux := http.NewServeMux()

	// Meter routes subsystem
	socketMux.HandleFunc("/meter", services.MeterHandler)

	// Protected routes subsystem
	protectedMux := http.NewServeMux()

	// General purpose dynamic dashboard
	protectedMux.HandleFunc("/dashboard", services.Dashboard)
	// Register new routes for account and balance pages for consumers:
	protectedMux.HandleFunc("/account", services.ConsumerAccount)
	protectedMux.HandleFunc("/balance", services.ConsumerBalance)

	protectedMux.HandleFunc("/logout", services.Logout)

	// Apply auth-specific middleware to protected routes
	protectedHandler := middleware.ChainMiddleware(
		middleware.WebAuthMiddleware,
	)(protectedMux)

	// Mount protected routes under /client/ path
	baseMux.Handle("/client/", http.StripPrefix("/client", protectedHandler))

	// Apply auth-specific middleware to websocket routes
	socketdHandler := middleware.ChainMiddleware(
	// middleware.AuthMiddleware,
	// REMINDER: add middleware for the websockets  here
	)(socketMux)

	// Mount WebSocket routes under /socket/ path
	baseMux.Handle("/v1/", http.StripPrefix("/v1", socketdHandler))

	// Apply GLOBAL middleware to ALL routes
	globalHandler := middleware.ChainMiddleware(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware, // Example of additional global middleware
	)(baseMux)

	return globalHandler
}
