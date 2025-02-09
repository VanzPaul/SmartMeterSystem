package services

import (
	"context"
	"fmt"
	"time"

	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDocument(db controllers.Database, collName models.Collection, data interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate the data before inserting
	if err := utils.ValidateData(data, collName); err != nil {
		return primitive.NilObjectID, fmt.Errorf("validation failed: %v", err)
	}

	dataBson, _ := bson.Marshal(data)

	// Insert the document
	insertResult, err := db.Create(ctx, dataBson)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %v", err)
	}

	// Ensure the inserted ID is of type primitive.ObjectID
	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("inserted ID is not of type primitive.ObjectID")
	}

	return insertedID, nil
}
