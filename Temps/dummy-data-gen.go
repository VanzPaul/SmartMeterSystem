package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type CommandType string

const (
	CommandTypeMeterRead CommandType = "METER_READ"
	CommandTypeReboot    CommandType = "REBOOT"
)

type CommandStatus string

const (
	CommandStatusPending   CommandStatus = "PENDING"
	CommandStatusCompleted CommandStatus = "COMPLETED"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// Core account structure
type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AccountNumber string             `bson:"accountNumber"`
	Name          string             `bson:"name"`
	Address       Address            `bson:"address"`
	Contact       Contact            `bson:"contact"`
	AccountType   AccountType        `bson:"accountType"`
	ConsumerType  ConsumerType       `bson:"consumerType"`
	Meters        []Meter            `bson:"meters"`
	Billing       Billing            `bson:"billing"`
	Ledger        Ledger             `bson:"ledger"`
	CreatedAt     int64              `bson:"createdAt"`
	UpdatedAt     int64              `bson:"updatedAt"`
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

type Meter struct {
	MeterID  string      `bson:"meterId"`
	Location GeoJSON     `bson:"location"`
	SIM      SIM         `bson:"sim"`
	Usage    []Usage     `bson:"usage"`
	Alerts   Alert       `bson:"alerts"`
	Commands Commands    `bson:"commands"`
	Status   MeterStatus `bson:"status"`
}

type GeoJSON struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type SIM struct {
	ICCID          string  `bson:"iccid"`
	MobileNumber   string  `bson:"mobileNumber"`
	DataUsage      float64 `bson:"dataUsage"`
	ActivationDate int64   `bson:"activationDate"`
}

type Usage struct {
	Date        int64   `bson:"date"`
	Consumption float64 `bson:"consumption"`
	Unit        string  `bson:"unit"`
}

type Alert struct {
	Current CurrentAlert `bson:"current"`
	History []AlertEvent `bson:"history"`
}

type CurrentAlert struct {
	Outage AlertStatus `bson:"outage"`
	Tamper AlertStatus `bson:"tamper"`
}

type AlertStatus struct {
	Active bool   `bson:"active"`
	Since  *int64 `bson:"since,omitempty"`
}

type AlertEvent struct {
	Type      string `bson:"type"`
	StartDate int64  `bson:"startDate"`
	EndDate   int64  `bson:"endDate"`
	Resolved  bool   `bson:"resolved"`
}

type Commands struct {
	Active  []ActiveCommand  `bson:"active"`
	History []HistoryCommand `bson:"history"`
}

type ActiveCommand struct {
	CommandID  string                 `bson:"commandId"`
	Type       CommandType            `bson:"type"`
	IssuedAt   int64                  `bson:"issuedAt"`
	Parameters map[string]interface{} `bson:"parameters"`
	Status     CommandStatus          `bson:"status"`
}

type HistoryCommand struct {
	ActiveCommand `bson:",inline"`
	CompletedAt   int64  `bson:"completedAt"`
	Response      string `bson:"response"`
}

type MeterStatus struct {
	LastSeen       int64   `bson:"lastSeen"`
	GridConnection bool    `bson:"gridConnection"`
	BatteryLevel   float64 `bson:"batteryLevel"`
}

type Billing struct {
	CurrentBill    CurrentBill `bson:"currentBill"`
	PaymentHistory []Payment   `bson:"paymentHistory"`
}

type CurrentBill struct {
	BillingPeriod    DateRange       `bson:"billingPeriod"`
	DueDate          int64           `bson:"dueDate"`
	TotalConsumption decimal.Decimal `bson:"totalConsumption"`
	AmountDue        decimal.Decimal `bson:"amountDue"`
	Paid             bool            `bson:"paid"`
}

type DateRange struct {
	Start int64 `bson:"start"`
	End   int64 `bson:"end"`
}

type Payment struct {
	BillingPeriod DateRange       `bson:"billingPeriod"`
	AmountPaid    decimal.Decimal `bson:"amountPaid"`
	PaymentDate   int64           `bson:"paymentDate"`
	PaymentMethod PaymentMethod   `bson:"paymentMethod"`
	TransactionID string          `bson:"transactionId"`
}

type Ledger struct {
	CurrentBalance Balance      `bson:"currentBalance"`
	UpcomingBill   UpcomingBill `bson:"upcomingBill"`
}

type Balance struct {
	Amount    decimal.Decimal `bson:"amount"`
	UpdatedAt int64           `bson:"updatedAt"`
}

type UpcomingBill struct {
	EstimatedAmount decimal.Decimal `bson:"estimatedAmount"`
	ProjectionDate  int64           `bson:"projectionDate"`
}

func main() {
	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://vanspaul09:ab7vSvvo14nx7gN3@cluster0.euhiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the database
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	// Get collection
	collection := client.Database("utilityDB").Collection("accounts")

	// Create sample account
	sampleAccount := Account{
		AccountNumber: "ACC-2023-001",
		Name:          "John Deere",
		Address: Address{
			Street:     "kaylaway",
			City:       "Nasugbu",
			State:      "Batangas",
			PostalCode: "4231",
		},
		Contact: Contact{
			Phone: "+639000000000",
			Email: "john.deere@example.com",
		},
		AccountType:  "consumer",
		ConsumerType: "residential",
		Meters: []Meter{
			{
				MeterID: "MTR-001",
				Location: GeoJSON{
					Type:        "Point",
					Coordinates: []float64{-73.935242, 40.730610},
				},
				SIM: SIM{
					ICCID:        "8910042348034555366",
					MobileNumber: "+639000000000",
					DataUsage:    10.0,
				},
				Usage: []Usage{
					{
						Date:        time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
						Consumption: 150.5,
						Unit:        "kWh",
					},
				},
				Alerts: Alert{
					Current: CurrentAlert{
						Outage: AlertStatus{Active: false},
						Tamper: AlertStatus{Active: false},
					},
					History: []AlertEvent{},
				},
				Commands: Commands{
					Active: []ActiveCommand{
						{
							CommandID: "CMD-001",
							Type:      "METER_READ",
							IssuedAt:  time.Now().Unix(),
							Parameters: map[string]interface{}{
								"interval": "15m",
							},
							Status: "PENDING",
						},
					},
					History: []HistoryCommand{},
				},
				Status: MeterStatus{
					LastSeen:       time.Now().Add(-30 * time.Minute).Unix(),
					GridConnection: true,
					BatteryLevel:   85.5,
				},
			},
		},
		Billing: Billing{
			CurrentBill: CurrentBill{
				BillingPeriod: DateRange{
					Start: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC).Unix(),
					End:   time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC).Unix(),
				},
				DueDate:          time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC).Unix(),
				TotalConsumption: decimal.NewFromFloat(450.0),
				AmountDue:        decimal.NewFromFloat(85.50),
				Paid:             false,
			},
			PaymentHistory: []Payment{
				{
					BillingPeriod: DateRange{
						Start: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC).Unix(),
						End:   time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC).Unix(),
					},
					AmountPaid:    decimal.NewFromFloat(82.75),
					PaymentDate:   time.Date(2023, 9, 14, 0, 0, 0, 0, time.UTC).Unix(),
					PaymentMethod: "credit_card",
					TransactionID: "TX-20230914-001",
				},
			},
		},
		Ledger: Ledger{
			CurrentBalance: Balance{
				Amount:    decimal.NewFromFloat(-85.50),
				UpdatedAt: time.Now().Unix(),
			},
			UpcomingBill: UpcomingBill{
				EstimatedAmount: decimal.NewFromFloat(90.00),
				ProjectionDate:  time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Insert document
	insertResult, err := collection.InsertOne(context.TODO(), sampleAccount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
}
