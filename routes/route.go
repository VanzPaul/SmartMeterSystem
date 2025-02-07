package routes

import (
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/middleware"
	"github.com/vanspaul/SmartMeterSystem/services"
)

func ServeMuxInit() http.Handler {
	// Register public routes
	mux := http.NewServeMux()

	// Register protected routes
	protectedMux := http.NewServeMux()

	// Group routes together
	mux.HandleFunc("POST /register", services.Register) // Register public routes
	mux.HandleFunc("POST /login", services.Login)
	mux.HandleFunc("/logout", services.Logout)
	protectedMux.HandleFunc("POST /dashboard", services.Dashboard) // Register protected routes

	// Apply AuthMiddleware
	protectedMiddleware := middleware.ChainMiddleware( // Protected Middlewares
		middleware.AuthMiddleware,
	)(protectedMux)

	globalMiddleware := middleware.ChainMiddleware( // General Middlewares
		middleware.LoggingMiddleware,
	)(mux)

	// Register the protected routes under a common prefix
	mux.Handle("/client/", http.StripPrefix("/client", protectedMiddleware))

	return globalMiddleware
}
