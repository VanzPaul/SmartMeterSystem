/*
 * @file internal/server/routes/v1.go
 * @brief v1.go file holds all v1 route groups and their handlers
 */
package routes

import (
	"SmartMeterSystem/cmd/web"
	"SmartMeterSystem/internal"
	"SmartMeterSystem/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// V1 Route Groups
type V1ConsumerRoute struct {
	Deps ServerDeps // Add this
}

// Other V1 route groups...
type V1MeterRoute struct {
	Deps ServerDeps
}

type V1EmployeeRoute struct {
	Deps ServerDeps
}

// V1Routes holds v1 route groups
type V1Routes struct {
	Consumer V1ConsumerRoute
	Meter    V1MeterRoute
	Employee V1EmployeeRoute
}

// V1Handler registers all v1 routes
func (r *V1Routes) V1Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.NotFound().Render(r.Context(), w)
	})

	// Register consumer routes
	mux.Handle("/consumer/", http.StripPrefix("/consumer", r.Consumer.HandleV1()))
	// Register employee routes
	mux.Handle("/employee/", http.StripPrefix("/employee", r.Employee.HandleV1()))
	return mux
}

// Handler methods for V1 route groups
func (c *V1ConsumerRoute) HandleV1() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		web.NotFound().Render(r.Context(), w)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// Inside your handler function
		userType := r.URL.Query().Get("user_type")

		web.LoginWebPage(c.Deps.GetDefaultRouteVersion(), userType).Render(r.Context(), w)
	})

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {

		web.ConsumerDashboardWebPage().Render(r.Context(), w)
	})

	return mux
}

func (c *V1EmployeeRoute) HandleV1() http.Handler {
	logger, loggerErr := internal.NewLogger()
	if loggerErr != nil {
		panic("Failed to create logger in v1.go")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		web.NotFound().Render(r.Context(), w)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userType := r.URL.Query().Get("user_type")
			web.LoginWebPage(c.Deps.GetDefaultRouteVersion(), userType).Render(r.Context(), w)
		case "POST":
			w.Header().Set("HX-Redirect", "/v1/employee/sysadmin/dashboard") // REMINDER: Make the redirect Dynamic
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	})

	// sysadminRoute := http.NewServeMux()
	sysadminRouteStruct := struct {
		dashboard struct {
			dashboard   http.HandlerFunc
			information http.HandlerFunc
		}
		consumer struct {
			consumer    http.HandlerFunc
			information http.HandlerFunc
		}
		accounts struct {
			accounts http.HandlerFunc
			forms    http.HandlerFunc
		}
		accounting struct {
			accounting http.HandlerFunc
			rates      http.HandlerFunc
		}
		payment struct {
			payment     http.HandlerFunc
			information http.HandlerFunc
		}
		logout http.HandlerFunc
	}{
		dashboard: struct {
			dashboard   http.HandlerFunc
			information http.HandlerFunc
		}{
			dashboard: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					web.SystemAdminEmployeeDashboardWebPage().Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			information: func(w http.ResponseWriter, r *http.Request) {
				// Extract the part after "/sysadmin/consumer/"
				pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/dashboard/")
				// Split to handle nested paths, take the first segment
				formType := strings.SplitN(pathPart, "/", 2)[0]

				switch r.Method {
				case "GET":
					switch formType {
					case "meter-list":
						smartmeters := []web.SmartMeter{
							{
								ID:        "SM001",
								Name:      "Smart Meter 1",
								Location:  "Calatagan",
								Latitude:  13.838432,
								Longitude: 120.632360,
								Status:    "active",
								Alert: []web.Alert{
									{
										ID:        "000000000",
										Type:      web.AlertTypePowerOutage,
										Timestamp: "1747312676",
										Status:    web.AlertStatusActive,
									},
								},
							},
							{
								ID:        "SM002",
								Name:      "Smart Meter 2",
								Location:  "Calatagan",
								Latitude:  13.839147,
								Longitude: 120.632257,
								Status:    "active",
							},
							{
								ID:        "SM003",
								Name:      "Meter 3",
								Location:  "Calatagan",
								Latitude:  13.838002,
								Longitude: 120.632220,
								Status:    "inactive",
								Alert: []web.Alert{
									{
										ID:        "000000001",
										Type:      web.AlertTypePowerOutage,
										Timestamp: "1747312676",
										Status:    web.AlertStatusActive,
									},
								},
							},
							{
								ID:        "SM002",
								Name:      "Meter 4",
								Location:  "Calatagan",
								Latitude:  13.837288,
								Longitude: 120.632164,
								Status:    "inactive",
							},
						}
						w.Header().Set("Content-Type", "application/json")
						if err := json.NewEncoder(w).Encode(smartmeters); err != nil {
							http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
							fmt.Println("Error:", err)
						}

					default:
						http.Error(w, "Not Found", http.StatusNotFound)
					}
				}
			},
		},
		consumer: struct {
			consumer    http.HandlerFunc
			information http.HandlerFunc
		}{
			consumer: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					web.SystemAdminEmployeeConsumerWebPage(
						[]web.ConsumerList{
							{
								ConsumerID:   "C001",
								ConsumerName: "John Doe",
								ConsumerType: web.ConsumerAccountTypeData.Residential,
								Status:       web.ConsumerAccountStatusData.Active,
							},
							{
								ConsumerID:   "C002",
								ConsumerName: "John Deer",
								ConsumerType: web.ConsumerAccountTypeData.Residential,
								Status:       web.ConsumerAccountStatusData.Inactive,
							},
						},
					).Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			information: func(w http.ResponseWriter, r *http.Request) {
				// Extract the part after "/sysadmin/consumer/"
				pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/consumer/")
				// Split to handle nested paths, take the first segment
				formType := strings.SplitN(pathPart, "/", 2)[0]

				switch r.Method {
				case "GET":
					switch formType {
					case "consumer-list":
						web.ConsumerListContainer(
							[]web.ConsumerList{
								{
									ConsumerID:   "C001",
									ConsumerName: "John Doe",
									ConsumerType: web.ConsumerAccountTypeData.Residential,
									Status:       web.ConsumerAccountStatusData.Active,
								},
								{
									ConsumerID:   "C002",
									ConsumerName: "John Deer",
									ConsumerType: web.ConsumerAccountTypeData.Residential,
									Status:       web.ConsumerAccountStatusData.Inactive,
								},
							},
						).Render(r.Context(), w)
					case "consumer-info":
						web.ConsumerInformationContainer().Render(r.Context(), w)
					// In your handler for "consumer-chart"
					case "consumer-chart":
						w.Header().Set("Content-Type", "text/html")

						// Lopgic Here

					}
				}
			},
		},
		accounts: struct {
			accounts http.HandlerFunc
			forms    http.HandlerFunc
		}{
			accounts: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					web.SystemAdminEmployeeAccountsWebPage().Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			forms: func(w http.ResponseWriter, r *http.Request) {
				// Extract the part after "/sysadmin/accounts/"
				pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/accounts/")
				// Split to handle nested paths, take the first segment
				formType := strings.SplitN(pathPart, "/", 2)[0]

				switch r.Method {
				case "GET":
					switch formType {
					case "meter-form":
						web.NewMeterAccountForm().Render(r.Context(), w)
					case "consumer-form":
						web.NewConsumerAccountForm().Render(r.Context(), w)
					case "employee-form":
						web.NewEmployeeAccountForm().Render(r.Context(), w)
					default:
						http.NotFound(w, r)
					}
				case "POST":
					switch formType {
					case "submit-meter-form":
						meterSN := r.FormValue("meter-sn")
						meterInstallationDate := r.FormValue("meter-installation-date")
						consumerTransformerId := r.FormValue("consumer-transformer-id")
						meterLatitude := r.FormValue("meter-latitude")
						meterLongitude := r.FormValue("meter-longitude")

						svc := database.New()

						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()

						// Create a new document
						insertResult, err := svc.InsertOne(ctx, "users", bson.M{
							"meter-sn":                meterSN,
							"meter-installation-date": meterInstallationDate, // Fixed field assignment
							"consumer-transformer-id": consumerTransformerId, // Fixed field assignment
							"meter-latitude":          meterLatitude,
							"meter-longitude":         meterLongitude,
						})

						if err != nil {
							logger.Sugar().Errorf("Insert failed: %v", err)
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return // ← Critical: Stop execution after error
						}

						if insertResult.InsertedID == nil {
							logger.Sugar().Warn("InsertedID is nil")
							http.Error(w, "No document created", http.StatusInternalServerError)
							return
						}

						logger.Sugar().Infof("Inserted ID: %v", insertResult.InsertedID)
						w.WriteHeader(http.StatusCreated)
					default:
						http.NotFound(w, r)
					}
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
		},
		accounting: struct {
			accounting http.HandlerFunc
			rates      http.HandlerFunc
		}{
			accounting: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					// web.SystemAdminEmployeeAccountingWebPage().Render(r.Context(), w)
					web.SystemAdminEmployeeAccountingWebPage(
						web.AccountingRatesTableFormType.Display,
						web.AccountingRatesTable{
							Date:        "01/01/01",
							Particulars: "RESIDENTIAL",
							Rates:       "",
							ERC:         "9.9298",
							AccountingRatesTableRowGroup: []web.AccountingRatesTableRowGroup{
								// Generation Charges
								{
									Particulars: "Generation Charges",
									Unit:        "",
									Rates:       "5.6092",
									ERC:         "5.6092",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Generation Energy Charge",
											Unit:        "PhP/kWh",
											Rates:       "5.6092",
											ERC:         "5.6092",
										},
										{
											Particulars: "Other Generation Rate Adjustment",
											Unit:        "PhP/kWh",
											Rates:       "0.0000",
											ERC:         "0.0000",
										},
									},
								},
								// Transmission Charges
								{
									Particulars: "Transmission Charges (NCCP)",
									Unit:        "",
									Rates:       "0.6853",
									ERC:         "0.6853",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Transmission Demand Charge",
											Unit:        "PhP/kW",
											Rates:       "0.0000",
											ERC:         "0.0000",
										},
										{
											Particulars: "Transmission System Charge",
											Unit:        "PhP/kWh",
											Rates:       "0.6853",
											ERC:         "0.6853",
										},
									},
								},
								// System Loss Charge
								{
									Particulars: "System Loss Charge",
									Unit:        "",
									Rates:       "0.9344",
									ERC:         "0.9344",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "System Loss Charge",
											Unit:        "PhP/kWh",
											Rates:       "0.9344",
											ERC:         "0.9344",
										},
									},
								},
								// Continue with other sections following the same pattern
								// Distribution Charges
								{
									Particulars: "Distribution Charges",
									Unit:        "",
									Rates:       "0.4613",
									ERC:         "0.4613",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Distribution Demand Charge",
											Unit:        "PhP/kW",
											Rates:       "0.0000",
											ERC:         "0.0000",
										},
										{
											Particulars: "Distribution System Charge",
											Unit:        "PhP/kWh",
											Rates:       "0.4613",
											ERC:         "0.4613",
										},
									},
								},
								// Supply Charges
								{
									Particulars: "Supply Charges",
									Unit:        "",
									Rates:       "0.5376",
									ERC:         "0.5376",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Supply Retail Customer Charge",
											Unit:        "PhP/Cust/Mo",
											Rates:       "0.0000",
											ERC:         "0.0000",
										},
										{
											Particulars: "Supply System Charge",
											Unit:        "PhP/kWh",
											Rates:       "0.5376",
											ERC:         "0.5376",
										},
									},
								},
								// Add remaining sections following the same structure...
								// Example for VAT section:
								{
									Particulars: "VAT",
									Unit:        "",
									Rates:       "1.0943",
									ERC:         "0.8543",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Generation",
											Unit:        "PhP/kWh",
											Rates:       "0.6376",
											ERC:         "0.6376",
										},
										{
											Particulars: "Transmission",
											Unit:        "PhP/kWh",
											Rates:       "0.1096",
											ERC:         "0.1096",
										},
										// Add other VAT components...
									},
								},
								// Universal Charge
								{
									Particulars: "Universal Charge",
									Unit:        "",
									Rates:       "0.2250",
									ERC:         "0.2250",
									SubRowGroup: []web.SubRowGroup{
										{
											Particulars: "Missionary Electrification",
											Unit:        "PhP/kWh",
											Rates:       "0.1822",
											ERC:         "0.1822",
										},
										// Add other universal charge components...
									},
								},
							},
						},
					).Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			rates: func(w http.ResponseWriter, r *http.Request) {
				// Extract the part after "/sysadmin/accounting/"
				pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/accounting/")
				// Split to handle nested paths, take the first segment
				formType := strings.SplitN(pathPart, "/", 2)[0]

				switch r.Method {
				case "GET":
					switch formType {
					case "update-rates-form":
						web.SystemAdminEmployeeAccountingTable(
							web.AccountingRatesTableFormType.FormRates,
							web.AccountingRatesTable{
								Date:        "01/01/01",
								Particulars: "RESIDENTIAL",
								Rates:       "",
								AccountingRatesTableRowGroup: []web.AccountingRatesTableRowGroup{
									// Generation Charges
									{
										Particulars: "Generation Charges",
										Unit:        "",
										Rates:       "5.6092",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Generation Energy Charge",
												Unit:        "PhP/kWh",
												Rates:       "5.6092",
											},
											{
												Particulars: "Other Generation Rate Adjustment",
												Unit:        "PhP/kWh",
												Rates:       "0.0000",
											},
										},
									},
									// Transmission Charges
									{
										Particulars: "Transmission Charges (NCCP)",
										Unit:        "",
										Rates:       "0.6853",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Transmission Demand Charge",
												Unit:        "PhP/kW",
												Rates:       "0.0000",
											},
											{
												Particulars: "Transmission System Charge",
												Unit:        "PhP/kWh",
												Rates:       "0.6853",
											},
										},
									},
									// System Loss Charge
									{
										Particulars: "System Loss Charge",
										Unit:        "",
										Rates:       "0.9344",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "System Loss Charge",
												Unit:        "PhP/kWh",
												Rates:       "0.9344",
											},
										},
									},
									// Continue with other sections following the same pattern
									// Distribution Charges
									{
										Particulars: "Distribution Charges",
										Unit:        "",
										Rates:       "0.4613",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Distribution Demand Charge",
												Unit:        "PhP/kW",
												Rates:       "0.0000",
											},
											{
												Particulars: "Distribution System Charge",
												Unit:        "PhP/kWh",
												Rates:       "0.4613",
											},
										},
									},
									// Supply Charges
									{
										Particulars: "Supply Charges",
										Unit:        "",
										Rates:       "0.5376",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Supply Retail Customer Charge",
												Unit:        "PhP/Cust/Mo",
												Rates:       "0.0000",
											},
											{
												Particulars: "Supply System Charge",
												Unit:        "PhP/kWh",
												Rates:       "0.5376",
											},
										},
									},
									// Add remaining sections following the same structure...
									// Example for VAT section:
									{
										Particulars: "VAT",
										Unit:        "",
										Rates:       "1.0943",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Generation",
												Unit:        "PhP/kWh",
												Rates:       "0.6376",
											},
											{
												Particulars: "Transmission",
												Unit:        "PhP/kWh",
												Rates:       "0.1096",
											},
											// Add other VAT components...
										},
									},
									// Universal Charge
									{
										Particulars: "Universal Charge",
										Unit:        "",
										Rates:       "0.2250",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Missionary Electrification",
												Unit:        "PhP/kWh",
												Rates:       "0.1822",
											},
											// Add other universal charge components...
										},
									},
								},
							},
						).Render(r.Context(), w)
					case "update-erc-form":
						web.SystemAdminEmployeeAccountingTable(
							web.AccountingRatesTableFormType.FormERC,
							web.AccountingRatesTable{
								Date:        "01/01/01",
								Particulars: "RESIDENTIAL",
								ERC:         "9.9298",
								AccountingRatesTableRowGroup: []web.AccountingRatesTableRowGroup{
									// Generation Charges
									{
										Particulars: "Generation Charges",
										Unit:        "",
										ERC:         "5.6092",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Generation Energy Charge",
												Unit:        "PhP/kWh",
												ERC:         "5.6092",
											},
											{
												Particulars: "Other Generation Rate Adjustment",
												Unit:        "PhP/kWh",
												ERC:         "0.0000",
											},
										},
									},
									// Transmission Charges
									{
										Particulars: "Transmission Charges (NCCP)",
										Unit:        "",
										ERC:         "0.6853",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Transmission Demand Charge",
												Unit:        "PhP/kW",
												ERC:         "0.0000",
											},
											{
												Particulars: "Transmission System Charge",
												Unit:        "PhP/kWh",
												ERC:         "0.6853",
											},
										},
									},
									// System Loss Charge
									{
										Particulars: "System Loss Charge",
										Unit:        "",
										ERC:         "0.9344",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "System Loss Charge",
												Unit:        "PhP/kWh",
												ERC:         "0.9344",
											},
										},
									},
									// Continue with other sections following the same pattern
									// Distribution Charges
									{
										Particulars: "Distribution Charges",
										Unit:        "",
										ERC:         "0.4613",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Distribution Demand Charge",
												Unit:        "PhP/kW",
												ERC:         "0.0000",
											},
											{
												Particulars: "Distribution System Charge",
												Unit:        "PhP/kWh",
												ERC:         "0.4613",
											},
										},
									},
									// Supply Charges
									{
										Particulars: "Supply Charges",
										Unit:        "",
										ERC:         "0.5376",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Supply Retail Customer Charge",
												Unit:        "PhP/Cust/Mo",
												ERC:         "0.0000",
											},
											{
												Particulars: "Supply System Charge",
												Unit:        "PhP/kWh",
												ERC:         "0.5376",
											},
										},
									},
									// Add remaining sections following the same structure...
									// Example for VAT section:
									{
										Particulars: "VAT",
										Unit:        "",
										ERC:         "0.8543",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Generation",
												Unit:        "PhP/kWh",
												ERC:         "0.6376",
											},
											{
												Particulars: "Transmission",
												Unit:        "PhP/kWh",
												ERC:         "0.1096",
											},
											// Add other VAT components...
										},
									},
									// Universal Charge
									{
										Particulars: "Universal Charge",
										Unit:        "",
										ERC:         "0.2250",
										SubRowGroup: []web.SubRowGroup{
											{
												Particulars: "Missionary Electrification",
												Unit:        "PhP/kWh",
												ERC:         "0.1822",
											},
											// Add other universal charge components...
										},
									},
								},
							},
						).Render(r.Context(), w)
					default:
						http.NotFound(w, r)
					}
				case "POST":
					switch formType {
					case "submit-update-rates-form":
						defer func() {
							if r := recover(); r != nil {
								http.Error(w, fmt.Sprintf("Internal server error: %v", r), http.StatusInternalServerError)
								log.Printf("Panic recovered: %v", r)
							}
						}()

						if r.Method != http.MethodPost {
							http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
							return
						}

						var payload web.AccountingRatesTable
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						// Process your data here...
						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// REMINDER: process submit-update-rates-form request here
						// w.Header().Set("Content-Type", "application/json")
						// if err := json.NewEncoder(w).Encode(payload); err != nil {
						// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
						// }

						// Process data here
						time.Sleep(1 * time.Second)
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("Success"))
					case "submit-update-erc-form":
						if r.Method != http.MethodPost {
							http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
							return
						}

						var payload web.AccountingRatesTable
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						// Process your data here...
						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// REMINDER: process submit-update-erc-form request here
						// w.Header().Set("Content-Type", "application/json")
						// if err := json.NewEncoder(w).Encode(payload); err != nil {
						// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
						// }

						time.Sleep(1 * time.Second)
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("Success"))
					default:
						http.NotFound(w, r)
					}
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
		},
		payment: struct {
			payment     http.HandlerFunc
			information http.HandlerFunc
		}{
			payment: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					web.SystemAdminEmployeePaymentWebPage().Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			information: func(w http.ResponseWriter, r *http.Request) {
				// Extract the part after "/sysadmin/consumer/"
				pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/payment/")
				// Split to handle nested paths, take the first segment
				formType := strings.SplitN(pathPart, "/", 2)[0]

				switch r.Method {
				case "GET":
					switch formType {
					// routes
					default:
						http.Error(w, "Not Found", http.StatusNotFound)
					}
				}
			},
		},
	}
	// System Admin Logout Route
	mux.HandleFunc("/sysadmin/logout", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("HX-Redirect", "/home")
		w.WriteHeader(http.StatusOK)
	})

	// System Admin Dashboard Routes
	mux.HandleFunc("/sysadmin/dashboard", sysadminRouteStruct.dashboard.dashboard)
	mux.HandleFunc("/sysadmin/dashboard/", sysadminRouteStruct.dashboard.information)
	// System Admin Consumer Routes
	mux.HandleFunc("/sysadmin/consumer", sysadminRouteStruct.consumer.consumer)
	mux.HandleFunc("/sysadmin/consumer/", sysadminRouteStruct.consumer.information)
	// System Admin Account Routes
	mux.HandleFunc("/sysadmin/accounts", sysadminRouteStruct.accounts.accounts)
	mux.HandleFunc("/sysadmin/accounts/", sysadminRouteStruct.accounts.forms)
	// System Admin Payment Routes
	mux.HandleFunc("/sysadmin/payment", sysadminRouteStruct.payment.payment)
	mux.HandleFunc("/sysadmin/payment/", sysadminRouteStruct.payment.information)
	// System Admin Accounting Routes
	mux.HandleFunc("/sysadmin/accounting", sysadminRouteStruct.accounting.accounting)
	mux.HandleFunc("/sysadmin/accounting/", sysadminRouteStruct.accounting.rates)

	return mux
}
