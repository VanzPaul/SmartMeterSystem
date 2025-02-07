package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMeter(uri, dbName, collName string, meter client.Meter) (primitive.ObjectID, error) {
	log.Print("Creating meter document")
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

	// Insert meter document
	dataBson, err := bson.Marshal(meter)
	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	log.Printf("Inserted document with ID: %v", insertResult.InsertedID)
	return insertedID, err
}

func CreateAccounting(uri, dbName, collName string, accounting client.Accounting) (primitive.ObjectID, error) {
	log.Print("Creating accounting document")
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

	// Insert accounting document
	dataBson, err := bson.Marshal(accounting)
	if err != nil {
		fmt.Println(err)
	}

	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
	return insertedID, err
}

func CreateConsumer(uri, dbName, collName string, consumer client.Consumer) (primitive.ObjectID, error) {
	log.Print("Creating consumer document")
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

	// Insert meter document
	dataBson, err := bson.Marshal(consumer)

	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
	return insertedID, err
}

func CreateAccount(uri, dbName, collName string, account models.Account) (primitive.ObjectID, error) {

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

	// Insert meter document
	dataBson, err := bson.Marshal(account)

	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
	return insertedID, err
}
