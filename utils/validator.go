package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/vanspaul/SmartMeterSystem/models"
	"go.uber.org/zap"
)

var validate *validator.Validate = validator.New()

// ValidateData validates data based on the collection name
func ValidateData(data interface{}, collName models.Collection) error {
	collectionStructMap := models.GetCcollectionMap()

	// Get the struct prototype for the given collection name
	structPrototype, exists := collectionStructMap[collName]
	if !exists {
		return fmt.Errorf("unknown collection: %s", collName)
	}

	// Log the types for debugging
	dataValue := reflect.ValueOf(data)
	prototypeValue := reflect.ValueOf(structPrototype)
	fmt.Printf("Data type: %v, Prototype type: %v\n", dataValue.Type(), prototypeValue.Type())

	// Ensure the data matches the expected struct type
	if dataValue.Type() != prototypeValue.Type() {
		Logger.Debug("data type mismatch", zap.String("collection", string(collName)), zap.Any("dataValueType", dataValue.Type()), zap.Any("prototypeValueType", prototypeValue.Type()))
		return fmt.Errorf("data type mismatch for collection %s", collName)
	}
	Logger.Debug("data type matched", zap.String("collection", string(collName)), zap.Any("dataValueType", dataValue.Type()), zap.Any("prototypeValueType", prototypeValue.Type()))

	// Perform validation using the validator package
	err := validate.Struct(data)
	if err != nil {
		// Handle validation errors
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Validation failed for field %s: %s\n", err.Field(), err.Tag())
		}
		return fmt.Errorf("validation failed: %v", err)
	}

	return nil
}
