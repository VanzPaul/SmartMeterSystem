package routes

import (
	"net/http"
)

func SetupRoutes() {
	// Group routes together
	http.HandleFunc("POST /register", Register)
	http.HandleFunc("POST /login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("POST /protected", Protected)
}
