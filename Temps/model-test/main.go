package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/users"
)

func main() {
	acc := models.Account{
		HashedPassword: "wpefojpanfasdfuivabnib",
		Email:          "johndoe@mail.com",
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		Role:           models.RoleConsumer,
		RoleSpecificData: users.Consumer{
			AccountNumber: "1234567890",
			Name:          "John Doe",
			Address: users.Address{
				Street:     "Calzada",
				City:       "Balayan",
				State:      "Batangas",
				PostalCode: "4230",
			},
			Contact: users.Contact{
				Phone: "09000000000",
				Email: "johndoe@mail.com",
			},
			Meters: []users.Meter{
				{
					MeterID: "adfbaludifbaufd",
					Location: users.GeoJSON{
						Type:        "point",
						Coordinates: []float64{12.34, 56.78},
					},
					SIM: users.SIM{
						ICCID:          "1830942872875982",
						MobileNumber:   "09000000000",
						DataUsageMb:    506985.0,
						ActivationDate: time.Now().Unix(),
					},
					Usage: []users.Usage{
						{
							Date:           time.Now().Unix(),
							ConsumptionKwh: 34456.0,
						},
					},
					Alerts: users.Alert{
						Current: users.CurrentAlert{
							Outage: users.AlertStatus{
								Active: false,
							},
							Tamper: users.AlertStatus{
								Active: false,
							},
						},
						History: []users.AlertEvent{
							{
								Type:      "outage",
								StartDate: time.Now().Unix(),
								EndDate:   time.Now().Unix(),
								Resolved:  true,
							},
						},
					},
					Commands: users.Commands{
						Active: []users.ActiveCommand{
							{
								CommandID:  "adhfoaidfaadf",
								Type:       users.CommandTypeMeterRead,
								IssuedAt:   time.Now().Unix(),
								Parameters: map[string]interface{}{},
								Status:     users.CommandStatusPending,
							},
						},
						History: []users.HistoryCommand{
							{
								CompletedAt: time.Now().Unix(),
								Response:    "Ok",
							},
						},
					},
					Status: users.MeterStatus{
						LastSeen:       time.Now().Unix(),
						GridConnection: true,
						BatteryLevel:   1.0,
					},
				},
			},
			Billing: users.Billing{
				CurrentBill: users.CurrentBill{
					BillingPeriod: users.DateRange{
						Start: time.Now().Unix(),
						End:   time.Now().Unix(),
					},
				},
				PaymentHistory: []users.Payment{
					{
						BillingPeriod: users.DateRange{
							Start: time.Now().Unix(),
							End:   time.Now().Unix(),
						},
						AmountPaid:    decimal.NewFromFloat(3545.05),
						PaymentDate:   time.Now().Unix(),
						PaymentMethod: users.PaymentMethodBankTransfer,
						TransactionID: "edafhkodouiafhoafh",
					},
				},
			},
			Ledger: users.Ledger{
				CurrentBalance: users.Balance{
					Amount:    decimal.NewFromFloat(4000.00),
					UpdatedAt: time.Now().Unix(),
				},
				UpcomingBill: users.UpcomingBill{
					EstimatedAmount: decimal.NewFromFloat(2000.00),
					ProjectionDate:  time.Now().Unix(),
				},
			},

			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		Status: models.AccountStatus{
			IsActive: true,
		},
	}

	jsonF, _ := json.MarshalIndent(acc, "", "	")
	fmt.Println(string(jsonF))
}
