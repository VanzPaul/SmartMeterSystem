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
	baseMux.HandleFunc("POST /register/submit", services.CreateGeneralAccount)

	baseMux.HandleFunc("/login", services.Login)
	baseMux.HandleFunc("POST /login/submit", services.SubmitWebLogin)

	// Protected routes subsystem
	socketMux := http.NewServeMux()

	// Meter routes subsystem
	socketMux.HandleFunc("/meter", services.MeterHandler)

	// Protected routes subsystem
	protectedMux := http.NewServeMux()

	// Routes for consumer
	protectedMux.HandleFunc("/consumer/dashboard", services.ConsumerDashboard)
	protectedMux.HandleFunc("/consumer/account", services.ConsumerAccount)
	protectedMux.HandleFunc("/consumer/balance", services.ConsumerBalance)

	// Routes for Field Administrator
	// REMINDER: replace with protected mux
	baseMux.HandleFunc("POST /fieldadmin/newmeter", services.CreateMeterAccount)

	// Routes for
	protectedMux.HandleFunc("/logout", services.Logout)

	// Apply auth-specific middleware to protected routes
	protectedHandler := middleware.ChainMiddleware(
		middleware.WebAuthMiddleware,
	)(protectedMux)

	// Mount protected routes under /client/ path
	baseMux.Handle("/client/", http.StripPrefix("/client", protectedHandler))

	/* // Apply auth-specific middleware to websocket routes
	socketdHandler := middleware.ChainMiddleware(
	middleware.AuthMiddleware,
	// TODO:: add middleware for mete websockets websockets  here
	)(socketMux)

	// Mount WebSocket routes under /socket/ path
	baseMux.Handle("/v1/", http.StripPrefix("/v1", socketdHandler)) */

	// Apply GLOBAL middleware to ALL routes
	globalHandler := middleware.ChainMiddleware(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware, // Example of additional global middleware
	)(baseMux)

	return globalHandler
}
