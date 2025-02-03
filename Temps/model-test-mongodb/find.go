package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vanspaul/SmartMeterSystem/controllers"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// MongoDB connection details
	uri := "mongodb+srv://vanspaul09:ab7vSvvo14nx7gN3@cluster0.euhiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	dbName := "test_db"
	collName := "persons"

	// Create a new MongoDB controller
	controller, err := controllers.NewMongoDBController(uri, dbName, collName)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := controller.Close(context.Background()); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Initialize the controller (optional if using NewMongoDBController)
	if err := controller.Init(context.Background()); err != nil {
		log.Fatalf("Failed to initialize MongoDB controller: %v", err)
	}

	// Find a single document
	var foundUser bson.M
	filter := bson.M{"email": "johndoe@mail.com"}
	if err := controller.FindOne(context.Background(), filter, &foundUser); err != nil {
		log.Fatalf("Failed to find document: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)
}
