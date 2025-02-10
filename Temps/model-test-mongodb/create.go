package main

import (
	"context"
	"log"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"github.com/vanspaul/SmartMeterSystem/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// REFERENCE: for creating documents
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load environment variables and initialize the logger
	if err := config.LoadEnv(); err != nil {
		// Use standard log if LoadEnv fails before initializing the logger
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	defer config.Logger.Sync()

	config.Logger.Debug("debug log") // This will now work if DEBUG=true

	// Rest of your code
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Create sample meter document
	// meter := client.Meter{
	// 	MeterID: "adfbaludifbaufd",
	// 	Location: client.GeoJSON{
	// 		Type:        "point",
	// 		Coordinates: []float64{12.34, 56.78},
	// 	},
	// 	SIM: client.SIM{
	// 		ICCID:          "1830942872875982",
	// 		MobileNumber:   "09000000000",
	// 		DataUsageMb:    506985.0,
	// 		ActivationDate: time.Now().Unix(),
	// 	},
	// 	Usage: []client.Usage{
	// 		{
	// 			Date: 1738490036,
	// 			Kwh:  5.2,
	// 		},
	// 		{
	// 			Date: 1738493636,
	// 			Kwh:  3.8,
	// 		},
	// 		{
	// 			Date: 1738497236,
	// 			Kwh:  7.1,
	// 		},
	// 	},
	// 	Alerts: client.Alert{
	// 		Current: client.CurrentAlert{
	// 			Outage: client.AlertStatus{
	// 				Active: false,
	// 			},
	// 			Tamper: client.AlertStatus{
	// 				Active: false,
	// 			},
	// 		},
	// 		History: []client.AlertEvent{
	// 			{
	// 				Type:      "outage",
	// 				StartDate: time.Now().Unix(),
	// 				EndDate:   time.Now().Unix(),
	// 				Resolved:  true,
	// 			},
	// 		},
	// 	},
	// 	Commands: client.Commands{
	// 		Active: []client.ActiveCommand{
	// 			{
	// 				CommandID:  "adhfoaidfaadf",
	// 				Type:       client.CommandTypeMeterRead,
	// 				IssuedAt:   time.Now().Unix(),
	// 				Parameters: map[string]interface{}{},
	// 				Status:     client.CommandStatusPending,
	// 			},
	// 		},
	// 		History: []client.HistoryCommand{
	// 			{
	// 				CompletedAt: time.Now().Unix(),
	// 				Response:    "Ok",
	// 			},
	// 		},
	// 	},
	// 	Status: client.MeterStatus{
	// 		LastSeen:       time.Now().Unix(),
	// 		GridConnection: true,
	// 		BatteryLevel:   1.0,
	// 	},
	// }

	/* 	accounting := client.Accounting{
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
	   	}
	*/

	meterid, _ := primitive.ObjectIDFromHex("67a9c1f786832b4ee56719e5")
	accountingId, _ := primitive.ObjectIDFromHex("67a9c60da9c810679b5eb5d5")

	var meterSlice []primitive.ObjectID
	meterSlice = append(meterSlice, meterid)

	// Insert meter document
	consumer := client.Consumer{
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
		MeterIDs:     meterSlice,
		AccountingID: accountingId,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	insertResult, createErr := services.CreateDocument(ctx, db, models.Consumers, &consumer)
	if createErr != nil {
		log.Fatalf("Err creating document %s: %v\n", models.Consumers, createErr)
	}
	log.Println("Created new Document")
	config.Logger.Info(insertResult.Hex())

	// Example usage (commented out for now)
	// var foundUser bson.M
	// filter := bson.M{"email": "johndoe@mail.com"}
	// if err := controllers.Database.FindOne(db, ctx, models.Collection.Consumers, filter, &foundUser); err != nil {
	// 	log.Fatalf("Failed to find document: %v", err)
	// }
	// fmt.Printf("Found user: %+v\n", foundUser)
}
