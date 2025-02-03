package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vanspaul/SmartMeterSystem/controllers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Implement this sample usage of /controllers/database.go
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

	// Example document
	user := bson.M{
		"name":  "John Doe",
		"email": "john.doe@example.com",
		"age":   30,
	}

	// Insert a document
	insertResult, err := controller.Create(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	// Find a single document
	var foundUser bson.M
	filter := bson.M{"email": "john.doe@example.com"}
	if err := controller.FindOne(context.Background(), filter, &foundUser); err != nil {
		log.Fatalf("Failed to find document: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	// Update a document
	update := bson.M{"age": 31}
	updateResult, err := controller.Update(context.Background(), filter, update)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}
	fmt.Printf("Updated %v document(s)\n", updateResult.ModifiedCount)

	// Replace a document
	newUser := bson.M{
		"name":  "Jane Doe",
		"email": "jane.doe@example.com",
		"age":   28,
	}
	replaceResult, err := controller.Replace(context.Background(), filter, newUser)
	if err != nil {
		log.Fatalf("Failed to replace document: %v", err)
	}
	fmt.Printf("Replaced %v document(s)\n", replaceResult.ModifiedCount)

	// Paginated query
	var users []bson.M
	page, perPage := int64(1), int64(10)
	if err := controller.PaginatedFind(context.Background(), bson.M{}, page, perPage, &users); err != nil {
		log.Fatalf("Failed to fetch paginated results: %v", err)
	}
	fmt.Printf("Paginated users: %+v\n", users)

	// Delete a document
	deleteResult, err := controller.Delete(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to delete document: %v", err)
	}
	fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)

	// Transaction example
	err = controller.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert a document within the transaction
		_, err := controller.Create(sessCtx, bson.M{
			"name":  "Transaction User",
			"email": "transaction.user@example.com",
			"age":   25,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to insert document in transaction: %w", err)
		}

		// Simulate an error to abort the transaction
		// Uncomment the line below to test transaction rollback
		// return nil, errors.New("simulated transaction error")

		return nil, nil
	})
	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}
}
