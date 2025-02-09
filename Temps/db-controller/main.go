package main

import (
	"log"
	"time"

	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"github.com/vanspaul/SmartMeterSystem/services"
)

// TODO: Implement this sample usage of /controllers/database.go
func main() {
	// MongoDB connection details
	uri := "mongodb+srv://vanspaul09:ab7vSvvo14nx7gN3@cluster0.euhiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	dbName := "system_data"

	// Create a new MongoDB controller
	db, dbErr := controllers.NewMongoDB(uri, dbName)
	if dbErr != nil {
		log.Printf("Err creating new MongoDB: %s\n", dbErr)
	}
	log.Println("Created new MongoDB")

	// Create Bson meter document
	meter := client.Meter{
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
	}

	meterId, createErr := services.CreateDocument(db, models.Meters, &meter)
	if createErr != nil {
		log.Fatalf("Err creating document %s: %v\n", models.Meters, createErr)
	}
	log.Println("Created new Document")

	log.Printf("Meter ID: %s", meterId)

}
