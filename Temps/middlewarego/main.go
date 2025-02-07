package main

import (
	"fmt"
	"net/http"
)

// AuthMiddleware checks for a valid Authorization header.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs incoming requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware chains multiple middleware functions together.
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// PublicHandler handles public routes.
func PublicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a public route.")
}

// ProtectedHandler handles protected routes.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected route.")
}

func main() {
	// Create a new ServeMux to register routes
	mux := http.NewServeMux()

	// Public route
	mux.HandleFunc("/public", PublicHandler)

	// Grouped protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/protected", ProtectedHandler)

	// Apply AuthMiddleware only to protected routes
	protectedHandlerWithAuth := ChainMiddleware(
		AuthMiddleware,
	)(protectedMux)

	// Register the protected routes under a common prefix
	mux.Handle("/api/", http.StripPrefix("/api", protectedHandlerWithAuth))

	// Apply LoggingMiddleware globally (to all routes)
	globalMiddleware := LoggingMiddleware(mux)

	// Start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", globalMiddleware)
}
