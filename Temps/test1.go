package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleOther(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a non domain request")
	w.Write([]byte("Hello, stranger..."))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the username and password from the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Print the username and password to the terminal
	log.Printf("Received login request - Username: %s, Password: %s\n", username, password)

	// Respond to the client
	w.Write([]byte("Login request received\n"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", handleOther)
	router.HandleFunc("/login", handleLogin)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
