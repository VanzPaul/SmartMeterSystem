package routes

import (
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/services"
)

func SetupRoutes() {
	// Group routes together
	http.HandleFunc("POST /register", services.Register)
	http.HandleFunc("POST /login", services.Login)
	http.HandleFunc("/logout", services.Logout)
	http.HandleFunc("POST /protected", services.Protected)
}
