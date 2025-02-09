package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/vanspaul/SmartMeterSystem/controllers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	var meterSlice []primitive.ObjectID
	meterid, metererr := createMeter()
	meterSlice = append(meterSlice, meterid)

	accountingid, accounterr := createAccounting()
	consumerid, consumererr := createConsumer(meterSlice, accountingid)
	accountid, accountingerr := createAccount(consumerid)
	// id, err := createMeter()
	if metererr != nil || accountingerr != nil || consumererr != nil || accounterr != nil {
		fmt.Println("err")
	}
	fmt.Printf("MeterId: %s\tAccountingId: %s\tConsumerId: %s\tAccountId: %s\n", meterid, accountingid, consumerid, accountid)
}

func createMeter() (primitive.ObjectID, error) {
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

	// Insert meter document
	dataBson, err := bson.Marshal(client.Meter{
		MeterID: "adfbaludifbaufd",
		Location: client.GeoJSON{
			Type:        "point",
			Coordinates: []float64{12.34, 56.78},
		},
		SIM: client.SIM{
			ICCID:          "1830942872875982",
			MobileNumber:   "09000000000",
			DataUsageMb:    506985.0,
			ActivationDate: time.Now().Unix(),
		},
		Usage: []client.Usage{
			{
				Date: 1738490036,
				Kwh:  5.2,
			},
			{
				Date: 1738493636,
				Kwh:  3.8,
			},
			{
				Date: 1738497236,
				Kwh:  7.1,
			},
		},
		Alerts: client.Alert{
			Current: client.CurrentAlert{
				Outage: client.AlertStatus{
					Active: false,
				},
				Tamper: client.AlertStatus{
					Active: false,
				},
			},
			History: []client.AlertEvent{
				{
					Type:      "outage",
					StartDate: time.Now().Unix(),
					EndDate:   time.Now().Unix(),
					Resolved:  true,
				},
			},
		},
		Commands: client.Commands{
			Active: []client.ActiveCommand{
				{
					CommandID:  "adhfoaidfaadf",
					Type:       client.CommandTypeMeterRead,
					IssuedAt:   time.Now().Unix(),
					Parameters: map[string]interface{}{},
					Status:     client.CommandStatusPending,
				},
			},
			History: []client.HistoryCommand{
				{
					CompletedAt: time.Now().Unix(),
					Response:    "Ok",
				},
			},
		},
		Status: client.MeterStatus{
			LastSeen:       time.Now().Unix(),
			GridConnection: true,
			BatteryLevel:   1.0,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	return insertedID, err
}

func createAccounting() (primitive.ObjectID, error) {
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

	// Insert meter document
	dataBson, err := bson.Marshal(client.Accounting{
		Billing: client.Billing{
			CurrentBill: client.CurrentBill{
				BillingPeriod: client.DateRange{
					Start: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC).Unix(),
					End:   time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC).Unix(),
				},
				DueDate:          time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC).Unix(),
				TotalConsumption: decimal.NewFromFloat(450.0),
				AmountDue:        decimal.NewFromFloat(85.50),
				Paid:             false,
			},
			PaymentHistory: []client.Payment{
				{
					BillingPeriod: client.DateRange{
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
		Ledger: client.Ledger{
			CurrentBalance: client.Balance{
				Amount:    decimal.NewFromFloat(-85.50),
				UpdatedAt: time.Now().Unix(),
			},
			UpcomingBill: client.UpcomingBill{
				EstimatedAmount: decimal.NewFromFloat(90.00),
				ProjectionDate:  time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	return insertedID, err
}

func createConsumer(meterId []primitive.ObjectID, accountingId primitive.ObjectID) (primitive.ObjectID, error) {
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

	// Insert meter document
	dataBson, err := bson.Marshal(client.Consumer{
		AccountNumber: "ACC-2023-001",
		Name:          "John Deere",
		Address: client.Address{
			Street:     "kaylaway",
			City:       "Nasugbu",
			State:      "Batangas",
			PostalCode: "4231",
		},
		Contact: client.Contact{
			Phone: "+639000000000",
			Email: "john.deere@example.com",
		},
		AccountType:  "consumer",
		ConsumerType: "residential",
		MeterIDs:     meterId,
		AccountingID: accountingId,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	})

	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	return insertedID, err
}

func createAccount(roledataid primitive.ObjectID) (primitive.ObjectID, error) {
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

	// Insert meter document
	dataBson, err := bson.Marshal(models.Account{
		HashedPassword:     "wpefojpanfasdfuivabnib",
		Email:              "johndoe@mail.com",
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		Role:               models.RoleConsumer,
		RoleSpecificDataID: roledataid,
	})

	if err != nil {
		fmt.Println(err)
	}
	insertResult, err := controller.Create(context.Background(), dataBson)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error: the value is not of type primitive.ObjectID
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}

	return insertedID, err
}
