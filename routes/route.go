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
	baseMux.HandleFunc("POST /register", services.Register)
	baseMux.HandleFunc("POST /login", services.Login)

	// Protected routes subsystem
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /dashboard", services.Dashboard)
	protectedMux.HandleFunc("POST /logout", services.Logout)

	// Apply auth-specific middleware to protected routes
	protectedHandler := middleware.ChainMiddleware(
		middleware.AuthMiddleware,
	)(protectedMux)

	// Mount protected routes under /client/ path
	baseMux.Handle("/client/", http.StripPrefix("/client", protectedHandler))

	// Apply GLOBAL middleware to ALL routes
	globalHandler := middleware.ChainMiddleware(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware, // Example of additional global middleware
	)(baseMux)

	return globalHandler
}
