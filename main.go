package main

import (
	"fmt"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/routes"
)

func main() {
	// Define all routes in a single function
	routes.SetupRoutes()

	// Start the server
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
