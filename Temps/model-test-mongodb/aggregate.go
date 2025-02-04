package main

import (
	"context"
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

	// Execute the aggregation query
	var foundUser bson.M
	cursor, err := controller.Aggregate(context.TODO(), pipeline, &foundUser)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
}
