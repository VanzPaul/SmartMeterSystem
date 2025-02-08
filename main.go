package main

import (
	"log"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/routes"
)

func main() {
	// Handle the routes
	globalMiddleware := routes.ServeMuxInit()

	// Start the server
	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", globalMiddleware)
}
