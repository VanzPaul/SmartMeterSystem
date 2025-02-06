package client

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Enums for type safety
type AccountType string

const (
	AccountTypeConsumer AccountType = "consumer"
	AccountTypeBusiness AccountType = "business"
)

type ConsumerType string

const (
	ConsumerTypeResidential ConsumerType = "residential"
	ConsumerTypeCommercial  ConsumerType = "commercial"
)

// Core account structure
type Consumer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AccountNumber string             `bson:"accountNumber"`
	Name          string             `bson:"name"`
	Address       Address            `bson:"address"`
	Contact       Contact            `bson:"contact"`
	AccountType   AccountType        `bson:"accountType"`
	ConsumerType  ConsumerType       `bson:"consumerType"`

	// Fields to store references to meters and accounting documents
	MeterIDs     []primitive.ObjectID `bson:"meterIds"`     // IDs of meter documents
	AccountingID primitive.ObjectID   `bson:"accountingId"` // IDs of accounting documents

	CreatedAt int64 `bson:"createdAt"`
	UpdatedAt int64 `bson:"updatedAt"`
}

type Address struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	State      string `bson:"state"`
	PostalCode string `bson:"postalCode"`
}

type Contact struct {
	Phone string `bson:"phone"`
	Email string `bson:"email"`
}
