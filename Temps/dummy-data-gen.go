package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AccountNumber string             `bson:"accountNumber"`
	Name          string             `bson:"name"`
	Address       Address            `bson:"address"`
	Contact       Contact            `bson:"contact"`
	AccountType   string             `bson:"accountType"`
	ConsumerType  string             `bson:"consumerType"`
	Meters        []Meter            `bson:"meters"`
	Billing       Billing            `bson:"billing"`
	Ledger        Ledger             `bson:"ledger"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

type Address struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	State      string `bson:"state"`
	PostalCode string `bson:"postalCode"`
	Country    string `bson:"country"`
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
	ICCID          string    `bson:"iccid"`
	IMSI           string    `bson:"imsi"`
	ActivationDate time.Time `bson:"activationDate"`
}

type Usage struct {
	Date        time.Time `bson:"date"`
	Consumption float64   `bson:"consumption"`
	Unit        string    `bson:"unit"`
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
	Active bool      `bson:"active"`
	Since  time.Time `bson:"since,omitempty"`
}

type AlertEvent struct {
	Type      string    `bson:"type"`
	StartDate time.Time `bson:"startDate"`
	EndDate   time.Time `bson:"endDate"`
	Resolved  bool      `bson:"resolved"`
}

type Commands struct {
	Active  []ActiveCommand  `bson:"active"`
	History []HistoryCommand `bson:"history"`
}

type ActiveCommand struct {
	CommandID  string                 `bson:"commandId"`
	Type       string                 `bson:"type"`
	IssuedAt   time.Time              `bson:"issuedAt"`
	Parameters map[string]interface{} `bson:"parameters"`
	Status     string                 `bson:"status"`
}

type HistoryCommand struct {
	ActiveCommand `bson:",inline"`
	CompletedAt   time.Time `bson:"completedAt"`
	Response      string    `bson:"response"`
}

type MeterStatus struct {
	LastSeen       time.Time `bson:"lastSeen"`
	GridConnection bool      `bson:"gridConnection"`
	BatteryLevel   float64   `bson:"batteryLevel"`
}

type Billing struct {
	CurrentBill    CurrentBill `bson:"currentBill"`
	PaymentHistory []Payment   `bson:"paymentHistory"`
}

type CurrentBill struct {
	BillingPeriod    DateRange `bson:"billingPeriod"`
	DueDate          time.Time `bson:"dueDate"`
	TotalConsumption float64   `bson:"totalConsumption"`
	AmountDue        float64   `bson:"amountDue"`
	Paid             bool      `bson:"paid"`
}

type DateRange struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}

type Payment struct {
	BillingPeriod DateRange `bson:"billingPeriod"`
	AmountPaid    float64   `bson:"amountPaid"`
	PaymentDate   time.Time `bson:"paymentDate"`
	PaymentMethod string    `bson:"paymentMethod"`
	TransactionID string    `bson:"transactionId"`
}

type Ledger struct {
	CurrentBalance Balance      `bson:"currentBalance"`
	UpcomingBill   UpcomingBill `bson:"upcomingBill"`
}

type Balance struct {
	Amount    float64   `bson:"amount"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type UpcomingBill struct {
	EstimatedAmount float64   `bson:"estimatedAmount"`
	ProjectionDate  time.Time `bson:"projectionDate"`
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
		Name:          "John Doe",
		Address: Address{
			Street:     "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
		Contact: Contact{
			Phone: "+1-555-123-4567",
			Email: "john.doe@example.com",
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
					ICCID:          "8910042348034555366",
					IMSI:           "310150123456789",
					ActivationDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				},
				Usage: []Usage{
					{
						Date:        time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
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
							IssuedAt:  time.Now(),
							Parameters: map[string]interface{}{
								"interval": "15m",
							},
							Status: "PENDING",
						},
					},
					History: []HistoryCommand{},
				},
				Status: MeterStatus{
					LastSeen:       time.Now().Add(-30 * time.Minute),
					GridConnection: true,
					BatteryLevel:   85.5,
				},
			},
		},
		Billing: Billing{
			CurrentBill: CurrentBill{
				BillingPeriod: DateRange{
					Start: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				DueDate:          time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
				TotalConsumption: 450.0,
				AmountDue:        85.50,
				Paid:             false,
			},
			PaymentHistory: []Payment{
				{
					BillingPeriod: DateRange{
						Start: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC),
					},
					AmountPaid:    82.75,
					PaymentDate:   time.Date(2023, 9, 14, 0, 0, 0, 0, time.UTC),
					PaymentMethod: "credit_card",
					TransactionID: "TX-20230914-001",
				},
			},
		},
		Ledger: Ledger{
			CurrentBalance: Balance{
				Amount:    -85.50,
				UpdatedAt: time.Now(),
			},
			UpcomingBill: UpcomingBill{
				EstimatedAmount: 90.00,
				ProjectionDate:  time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert document
	insertResult, err := collection.InsertOne(context.TODO(), sampleAccount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
}
