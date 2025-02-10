package middleware

import (
	"log"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.uber.org/zap"
)

// ChainMiddleware chains multiple middleware functions together.
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// AuthMiddleware validates the session token and CSRF token using the Authorize function.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the singleton store instance
		store := models.GetStore()

		// Call the Authorize function to validate the session and CSRF tokens
		if err := utils.Authorize(r, store.Users, store.Sessions); err != nil {
			utils.Logger.Error("Authentication failed", zap.Any("Err", err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs incoming requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Sugar().Infof("Request received: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// In middleware package
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
