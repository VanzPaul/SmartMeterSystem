package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vanspaul/SmartMeterSystem/config"
)

func main() {
	// Load environment variables from the .env file
	fmt.Println("Loading .env file...")
	config.LoadEnv()

	// Create a loader instance
	loader := config.DefaultEnvLoader{}

	// Connect to MongoDB
	fmt.Println("Connecting to MongoDB...")
	client := config.ConnectDB(loader)

	// Use the client for database operations
	defer func() {
		fmt.Println("Disconnecting from MongoDB...")
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
		fmt.Println("Disconnected from MongoDB.")
	}()
}
