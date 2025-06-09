/*
 * @file internal/server/routes/v1.go
 * @brief v1.go file holds all v1 route groups and their handlers
 */
package routes

import (
	"SmartMeterSystem/cmd/web"
	"SmartMeterSystem/internal"
	"SmartMeterSystem/internal/database"
	"SmartMeterSystem/internal/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	logger, loggerErr := internal.NewLogger()
	if loggerErr != nil {
		panic("Failed to create logger in v1.go")
	}
	logger.Info("")

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		userType := r.URL.Query().Get("user_type")

		switch r.Method {
		case "GET":
			web.LoginWebPage(c.Deps.GetDefaultRouteVersion(), userType).Render(r.Context(), w)
		case "POST":
			w.Header().Set("HX-Redirect", "/v1/consumer/dashboard") // REMINDER: Make the redirect Dynamic
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	})

	consumerRouteStruct := struct {
		dashboard struct {
			dashboard   http.HandlerFunc
			information http.HandlerFunc
		}
	}{
		dashboard: struct {
			dashboard   http.HandlerFunc
			information http.HandlerFunc
		}{
			dashboard: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					web.ConsumerPage().Render(r.Context(), w)
				default:
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
			},
			information: func(w http.ResponseWriter, r *http.Request) {
				// // Extract the part after "/sysadmin/consumer/"
				// pathPart := strings.TrimPrefix(r.URL.Path, "/sysadmin/dashboard/")
				// // Split to handle nested paths, take the first segment
				// formType := strings.SplitN(pathPart, "/", 2)[0]
			},
		},
	}

	mux.HandleFunc("/dashboard", consumerRouteStruct.dashboard.dashboard)
	mux.HandleFunc("/dashboard/", consumerRouteStruct.dashboard.information)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		web.NotFound().Render(r.Context(), w)
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
		userType := r.URL.Query().Get("user_type")

		switch r.Method {
		case "GET":
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
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()

						svc := database.New()

						cursorFind, errFindMany := svc.FindMany(ctx, "meters", bson.M{"isActive": true})
						if errFindMany != nil {
							c.Deps.GetLogger().Sugar().Errorf("Error fetching active meters: %v", errFindMany)
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return
						}
						defer cursorFind.Close(ctx)

						var smartmeters []bson.M
						if err := cursorFind.All(ctx, &smartmeters); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Error decoding cursor: %v", err)
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return
						}

						w.Header().Set("Content-Type", "application/json")
						if err := json.NewEncoder(w).Encode(smartmeters); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Error encoding JSON: %v", err)
							http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
							return
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
				// Split to handle nested paths
				pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/sysadmin/accounts/"), "/")

				switch r.Method {
				case "GET":
					// Check path length
					if len(pathParts) < 1 {
						http.NotFound(w, r)
						return
					}

					switch subpath := pathParts[0]; subpath {
					case "create":
						// Check path length
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}
						usecase := subpath
						switch subpath := pathParts[1]; subpath {
						case "meter-form":
							web.NewMeterAccountForm(usecase).Render(r.Context(), w)
						case "consumer-form":
							web.NewConsumerAccountForm(usecase).Render(r.Context(), w)
						case "employee-form":
							web.NewEmployeeAccountForm().Render(r.Context(), w)
						default:
							http.NotFound(w, r)
						}

					case "update":
						// Check path length
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}
						usecase := subpath
						switch subpath := pathParts[1]; subpath {
						case "meter-form":
							web.NewMeterAccountForm(usecase).Render(r.Context(), w)
						case "consumer-form":
							web.NewConsumerAccountForm(usecase).Render(r.Context(), w)
						case "employee-form":
							web.NewEmployeeAccountForm().Render(r.Context(), w)
						default:
							http.NotFound(w, r)
						}

					case "delete":
						// Check path length
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}
						usecase := subpath
						switch subpath := pathParts[1]; subpath {
						case "meter-form":
							web.NewMeterAccountForm(usecase).Render(r.Context(), w)
						case "consumer-form":
							web.NewConsumerAccountForm(usecase).Render(r.Context(), w)
						case "employee-form":
							web.NewEmployeeAccountForm().Render(r.Context(), w)
						default:
							http.NotFound(w, r)
						}

					default:
						http.NotFound(w, r)
					}
				case "POST":
					// Check path length
					if len(pathParts) < 1 {
						http.NotFound(w, r)
						return
					}

					switch subpath := pathParts[0]; subpath {
					case "create":
						// Check path length
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}

						switch subpath := pathParts[1]; subpath {
						case "create-meter-form":
							// Parse form values
							form_meterNo := r.FormValue("create-meter-no")
							form_meterInstallationDate := r.FormValue("create-meter-installation-date")
							form_consumerTransformerId := r.FormValue("create-consumer-transformer-id")
							form_meterLatitude := r.FormValue("create-meter-latitude")
							form_meterLongitude := r.FormValue("create-meter-longitude")
							form_consumerAccNo := r.FormValue("create-consumer-acc-no")
							form_meterAddress := r.FormValue("create-meter-address")

							// Validate required fields
							if form_meterNo == "" || form_meterInstallationDate == "" || form_consumerTransformerId == "" ||
								form_meterLatitude == "" || form_meterLongitude == "" || form_consumerAccNo == "" ||
								form_meterAddress == "" {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "All fields are required",
								})
								return
							}

							// Convert form values to appropriate types
							meterNo, err_meterNo := strconv.Atoi(form_meterNo)
							meterInstallationDate, err_meterInstallationDate := time.Parse("2006-01-02", form_meterInstallationDate)
							consumerTransformerId := form_consumerTransformerId
							meterLatitude, err_meterLatitude := strconv.ParseFloat(form_meterLatitude, 64)
							meterLongitude, err_meterLongitude := strconv.ParseFloat(form_meterLongitude, 64)
							consumerAccNo, err_consumerAccNo := strconv.Atoi(form_consumerAccNo)

							// Handle individual conversion errors
							if err_meterNo != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Meter Number",
								})
								return
							}

							if err_meterInstallationDate != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Installation Date format. Expected YYYY-MM-DD",
								})
								return
							}

							if err_meterLatitude != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Latitude value",
								})
								return
							}

							if err_meterLongitude != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Longitude value",
								})
								return
							}

							if err_consumerAccNo != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Consumer Account Number",
								})
								return
							}

							svc := database.New()

							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							bsonData, err_bson := bson.Marshal(models.MeterDocument{
								MeterNumber:           meterNo,
								ConsumerAccNo:         consumerAccNo,
								InstallationDate:      meterInstallationDate,
								ConsumerTransformerID: consumerTransformerId,
								Coordinates:           []float64{meterLongitude, meterLatitude},
								Address:               form_meterAddress,
								IsActive:              true,
								SmartMeter: models.SmartMeter{
									IsActive:              false,
									Alert:                 nil,
									UsageKwh:              79.00,
									ReadingHistory30days:  nil,
									ReadingHistory24hours: nil,
								},
							})

							if err_bson != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Error Parsing Meter Document",
								})
								return
							}

							// Create a new document
							insertResult, err := svc.InsertOne(ctx, "meters", bsonData)

							if err != nil {
								// Check for duplicate key error
								var writeException mongo.WriteException
								if errors.As(err, &writeException) {
									for _, writeError := range writeException.WriteErrors {
										if writeError.Code == 11000 { // MongoDB duplicate key error code
											w.Header().Set("Content-Type", "application/json")
											w.WriteHeader(http.StatusConflict)
											json.NewEncoder(w).Encode(map[string]string{
												"error": "Meter Number already exists",
											})
											return
										}
									}
								}

								// Other errors
								logger.Sugar().Errorf("Insert failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if insertResult.InsertedID == nil {
								logger.Sugar().Error("InsertedID is nil")
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							logger.Sugar().Infof("Inserted ID: %v", insertResult.InsertedID)
							w.WriteHeader(http.StatusCreated)

						case "create-consumer-form":
							form_accountNumber := r.FormValue("create-consumer-acc-no")
							form_firstName := r.FormValue("create-consumer-first-name")
							form_middleName := r.FormValue("create-consumer-middle-name")
							form_lastName := r.FormValue("create-consumer-last-name")
							form_suffix := r.FormValue("create-consumer-suffix-name")
							form_birthDate := r.FormValue("create-consumer-birth-date") // Note: Inconsistent field name
							form_province := r.FormValue("create-consumer-province")
							form_postalCode := r.FormValue("create-consumer-postal-code")
							form_cityMunicipality := r.FormValue("create-consumer-city-municipality")
							form_barangay := r.FormValue("create-consumer-barangay")
							form_street := r.FormValue("create-consumer-street")
							form_phoneNumber := r.FormValue("create-consumer-phone-number")

							// Validate required fields
							if form_accountNumber == "" || form_firstName == "" || form_middleName == "" || form_lastName == "" ||
								form_suffix == "" || form_birthDate == "" || form_province == "" || form_postalCode == "" ||
								form_cityMunicipality == "" || form_barangay == "" || form_street == "" || form_phoneNumber == "" {

								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "All fields are required",
								})
								return
							}

							// Convert form values to appropriate types
							accountNumber, errAccountNumber := strconv.Atoi(form_accountNumber)
							firstName := form_firstName
							middleName := form_middleName
							lastName := form_lastName
							suffix := form_suffix
							birthDate, errBirthDate := time.Parse("2006-01-02", form_birthDate)
							province := form_province
							postalCode, errPostalCode := strconv.Atoi(form_postalCode)
							cityMunicipality := form_cityMunicipality
							barangay := form_barangay
							street := form_street
							phoneNumber, errPhoneNumber := strconv.Atoi(form_phoneNumber)

							// Handle individual conversion errors

							if errAccountNumber != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Account Number",
								})
								return
							}

							if errBirthDate != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Birth Date format. Expected YYYY-MM-DD",
								})
								return
							}

							if errPostalCode != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Postal Code",
								})
								return
							}

							if errPhoneNumber != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Phone Number",
								})
								return
							}

							svc := database.New()
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							bsonConsumer, err_cons_bson := bson.Marshal(models.ConsumerDocument{
								ID:               accountNumber,
								AccountNumber:    accountNumber,
								FirstName:        firstName,
								MiddleName:       middleName,
								LastName:         lastName,
								Suffix:           suffix,
								BirthDate:        birthDate,
								Province:         province,
								PostalCode:       postalCode,
								CityMunicipality: cityMunicipality,
								Barangay:         barangay,
								Street:           street,
								PhoneNumber:      phoneNumber,
								IsActive:         true, // REMINDER: implement this in the frontend
							})
							bsonBalance, err_balance_bson := bson.Marshal(models.ConsumerBalanceDocument{
								ID:            accountNumber,
								AccountNumber: accountNumber,
								ConsumerType:  "RESIDENTIAL",
								IsActive:      true,
							})

							if err_cons_bson != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Error Parsing Consumer Document",
								})
								return
							} else if err_balance_bson != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Error Parsing Balance Document",
								})
								return
							}

							// Create a new document
							insertConsumerResult, err_consumer := svc.InsertOne(ctx, "consumers", bsonConsumer)
							insertBalanceResult, _ := svc.InsertOne(ctx, "balances", bsonBalance)

							if err_consumer != nil {
								// Check for duplicate key error
								var writeException mongo.WriteException
								if errors.As(err_consumer, &writeException) {
									for _, writeError := range writeException.WriteErrors {
										if writeError.Code == 11000 { // MongoDB duplicate key error code
											w.Header().Set("Content-Type", "application/json")
											w.WriteHeader(http.StatusConflict)
											json.NewEncoder(w).Encode(map[string]string{
												"error": "Consumer Account Number already exists",
											})
											return
										}
									}
								}

								// Other errors
								logger.Sugar().Errorf("Insert failed: %v", err_consumer)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if insertConsumerResult.InsertedID == nil {
								logger.Sugar().Error("Consumer InsertedID is nil")
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							} else if insertBalanceResult.InsertedID == nil {
								logger.Sugar().Error("Balance InsertedID is nil")
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							logger.Sugar().Infof("Inserted ID: %v", insertConsumerResult.InsertedID)
							w.WriteHeader(http.StatusCreated)

						default:
							http.NotFound(w, r)
						}

					case "update":
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}

						// FIXME: update this code according to the new data structure

						switch subpath := pathParts[1]; subpath {
						case "update-meter-form":
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							svc := database.New()

							// Verification endpoint
							if len(pathParts) > 2 && pathParts[2] == "verify" {
								form_meterNo := r.FormValue("meter-no")
								meterNo, err_meterNo := strconv.Atoi(form_meterNo)

								// Handle conversion error
								if err_meterNo != nil {
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusBadRequest)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Invalid Meter Number",
									})
									return
								}

								var existingMeter models.MeterDocument
								err := svc.FindOne(ctx, "meters", bson.M{"_id": meterNo}).Decode(&existingMeter)
								if err != nil {
									if err == mongo.ErrNoDocuments {
										logger.Sugar().Warnf("Meter not found: %s", meterNo)
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusNotFound)
										json.NewEncoder(w).Encode(map[string]string{
											"error": "Meter not found: " + form_meterNo,
										})
										return
									}
									logger.Sugar().Errorf("Database error: %v", err)
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusInternalServerError)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Internal server error",
									})
									return
								}

								logger.Sugar().Infof("Found existing meter: %v", existingMeter)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(map[string]interface{}{
									"meterNo":          existingMeter.ID,
									"consumerAccNo":    existingMeter.ConsumerAccNo,
									"installationDate": existingMeter.InstallationDate.Format("2006-01-02"),
									"transformerId":    existingMeter.ConsumerTransformerID,
									"coordinates":      existingMeter.Coordinates,
								})
								return
							}

							// Update handler
							meterNo := r.FormValue("meter-no")
							installationDateStr := r.FormValue("meter-installation-date")
							transformerID := r.FormValue("consumer-transformer-id")
							latitude := r.FormValue("meter-latitude")
							longitude := r.FormValue("meter-longitude")
							consumerAccNoStr := r.FormValue("consumer-acc-no")

							// Parse meter number
							meterNoInt, err := strconv.Atoi(meterNo)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid meter account number format",
								})
								return
							}
							// Parse installation date
							installationDate, err := time.Parse("2006-01-02", installationDateStr)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid date format, use YYYY-MM-DD",
								})
								return
							}

							// Parse consumer account number
							consumerAccNo, err := strconv.Atoi(consumerAccNoStr)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid consumer account number format",
								})
								return
							}

							// Parse latitude/longitude
							lat, err := strconv.ParseFloat(latitude, 64)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid latitude format",
								})
								return
							}

							lng, err := strconv.ParseFloat(longitude, 64)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid longitude format",
								})
								return
							}

							updateData := bson.M{
								"$set": bson.M{
									"installDate":   installationDate,
									"transformerId": transformerID,
									"lat":           lat,
									"long":          lng,
									"acctNo":        consumerAccNo,
								},
							}

							updateResult, err := svc.UpdateOne(ctx, "meters", bson.M{"_id": meterNoInt}, updateData)
							if err != nil {
								logger.Sugar().Errorf("Update failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if updateResult.MatchedCount == 0 {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusNotFound)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Meter no longer exists: " + meterNo,
								})
								return
							}

							logger.Sugar().Infof("Updated %d document(s)", updateResult.ModifiedCount)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success": true,
								"message": "Meter updated successfully",
							})

						case "update-consumer-form":
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							svc := database.New()

							// Verification endpoint
							if len(pathParts) > 2 && pathParts[2] == "verify" {
								form_consumerAccNo := r.FormValue("update-consumer-acc-no")
								accNo, err_accNo := strconv.Atoi(form_consumerAccNo)

								// Handle conversion error
								if err_accNo != nil {
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusBadRequest)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Invalid Account Number",
									})
									return
								}

								var existingConsumer models.ConsumerDocument
								err := svc.FindOne(ctx, "consumers", bson.M{"_id": accNo}).Decode(&existingConsumer)
								if err != nil {
									if err == mongo.ErrNoDocuments {
										logger.Sugar().Warnf("Meter not found: %s", form_consumerAccNo)
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusNotFound)
										json.NewEncoder(w).Encode(map[string]string{
											"error": "Meter not found: " + form_consumerAccNo,
										})
										return
									}
									logger.Sugar().Errorf("Database error: %v", err)
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusInternalServerError)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Internal server error",
									})
									return
								}

								logger.Sugar().Infof("Found existing meter: %v", existingConsumer)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(map[string]interface{}{
									"acctNum":    existingConsumer.AccountNumber,
									"firstName":  existingConsumer.FirstName,
									"middleName": existingConsumer.MiddleName,
									"lastName":   existingConsumer.LastName,
									"suffix":     existingConsumer.Suffix,
									"birthDate":  existingConsumer.BirthDate.Format("2006-01-02"),
									"province":   existingConsumer.Province,
									"postalCode": existingConsumer.PostalCode,
									"cityMun":    existingConsumer.CityMunicipality,
									"barangay":   existingConsumer.Barangay,
									"street":     existingConsumer.Street,
									"phoneNum":   existingConsumer.PhoneNumber,
								})
								return
							}

							// Update Handler
							form_accountNumber := r.FormValue("update-consumer-acc-no")
							form_firstName := r.FormValue("update-consumer-first-name")
							form_middleName := r.FormValue("update-consumer-middle-name")
							form_lastName := r.FormValue("update-consumer-last-name")
							form_suffix := r.FormValue("update-consumer-suffix-name")
							form_birthDate := r.FormValue("update-consumer-birth-date") // Note: Inconsistent field name
							form_province := r.FormValue("update-consumer-province")
							form_postalCode := r.FormValue("update-consumer-postal-code")
							form_cityMunicipality := r.FormValue("update-consumer-city-municipality")
							form_barangay := r.FormValue("update-consumer-barangay")
							form_street := r.FormValue("update-consumer-street")
							form_phoneNumber := r.FormValue("update-consumer-phone-number")

							// Validate required fields
							if form_accountNumber == "" || form_firstName == "" || form_middleName == "" || form_lastName == "" ||
								form_suffix == "" || form_birthDate == "" || form_province == "" || form_postalCode == "" ||
								form_cityMunicipality == "" || form_barangay == "" || form_street == "" || form_phoneNumber == "" {

								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "All fields are required",
								})
								return
							}

							// Convert form values to appropriate types
							accountNumber, errAccountNumber := strconv.Atoi(form_accountNumber)
							firstName := form_firstName
							middleName := form_middleName
							lastName := form_lastName
							suffix := form_suffix
							birthDate, errBirthDate := time.Parse("2006-01-02", form_birthDate)
							province := form_province
							postalCode, errPostalCode := strconv.Atoi(form_postalCode)
							cityMunicipality := form_cityMunicipality
							barangay := form_barangay
							street := form_street
							phoneNumber, errPhoneNumber := strconv.Atoi(form_phoneNumber)

							// Handle individual conversion errors

							if errAccountNumber != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Account Number",
								})
								return
							}

							if errBirthDate != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Birth Date format. Expected YYYY-MM-DD",
								})
								return
							}

							if errPostalCode != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Postal Code",
								})
								return
							}

							if errPhoneNumber != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Phone Number",
								})
								return
							}

							updateData := bson.M{
								"$set": bson.M{
									"acctNum":    accountNumber,
									"firstName":  firstName,
									"middleName": middleName,
									"lastName":   lastName,
									"suffix":     suffix,
									"birthDate":  birthDate,
									"province":   province,
									"postalCode": postalCode,
									"cityMun":    cityMunicipality,
									"barangay":   barangay,
									"street":     street,
									"phoneNum":   phoneNumber,
								},
							}

							updateResult, err := svc.UpdateOne(ctx, "consumers", bson.M{"_id": accountNumber}, updateData)
							if err != nil {
								logger.Sugar().Errorf("Update failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if updateResult.MatchedCount == 0 {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusNotFound)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Meter no longer exists: " + form_accountNumber,
								})
								return
							}

							logger.Sugar().Infof("Updated %d document(s)", updateResult.ModifiedCount)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success": true,
								"message": "Meter updated successfully",
							})

						default:
							http.NotFound(w, r)
						}

					case "delete":
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}

						switch subpath := pathParts[1]; subpath {
						case "delete-meter-form":
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							svc := database.New()

							// Verification endpoint
							if len(pathParts) > 2 && pathParts[2] == "verify" {
								form_meterNo := r.FormValue("meter-no")
								meterNo, err_meterNo := strconv.Atoi(form_meterNo)

								// Handle conversion error
								if err_meterNo != nil {
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusBadRequest)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Invalid Meter Number",
									})
									return
								}

								var existingMeter models.MeterDocument
								err := svc.FindOne(ctx, "meters", bson.M{"_id": meterNo}).Decode(&existingMeter)
								if err != nil {
									if err == mongo.ErrNoDocuments {
										logger.Sugar().Warnf("Meter not found: %s", form_meterNo)
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusNotFound)
										json.NewEncoder(w).Encode(map[string]string{
											"error": "Meter not found: " + form_meterNo,
										})
										return
									}
									logger.Sugar().Errorf("Database error: %v", err)
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusInternalServerError)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Internal server error",
									})
									return
								}

								logger.Sugar().Infof("Found existing meter: %v", existingMeter)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(map[string]interface{}{
									"meterNo":          existingMeter.ID,
									"consumerAccNo":    existingMeter.ConsumerAccNo,
									"installationDate": existingMeter.InstallationDate.Format("2006-01-02"),
									"transformerId":    existingMeter.ConsumerTransformerID,
									"coordinates":      existingMeter.Coordinates,
								})
								return
							}

							// Delete handler
							meterNo := r.FormValue("meter-no")

							// Parse meter number
							meterNoInt, err := strconv.Atoi(meterNo)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid meter account number format",
								})
								return
							}

							// Delete document using meter number as ID
							deleteResult, err := svc.DeleteOne(ctx, "meters", bson.M{"_id": meterNoInt})
							if err != nil {
								logger.Sugar().Errorf("Delete failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if deleteResult.DeletedCount == 0 {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusNotFound)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Meter not found: " + meterNo,
								})
								return
							}

							logger.Sugar().Infof("Deleted %d document(s)", deleteResult.DeletedCount)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success": true,
								"message": "Meter deleted successfully",
							})

						case "delete-consumer-form":
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							svc := database.New()

							// Verification endpoint
							if len(pathParts) > 2 && pathParts[2] == "verify" {
								form_accNo := r.FormValue("acc-no")
								accNo, err_accNo := strconv.Atoi(form_accNo)

								// Handle conversion error
								if err_accNo != nil {
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusBadRequest)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Invalid Account Number",
									})
									return
								}

								var existingConsumer models.ConsumerDocument
								err := svc.FindOne(ctx, "consumers", bson.M{"_id": accNo}).Decode(&existingConsumer)
								if err != nil {
									if err == mongo.ErrNoDocuments {
										logger.Sugar().Warnf("Consumer not found: %s", form_accNo)
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusNotFound)
										json.NewEncoder(w).Encode(map[string]string{
											"error": "Consumer Account not found: " + form_accNo,
										})
										return
									}
									logger.Sugar().Errorf("Database error: %v", err)
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusInternalServerError)
									json.NewEncoder(w).Encode(map[string]string{
										"error": "Internal server error",
									})
									return
								}

								logger.Sugar().Infof("Found existing meter: %v", existingConsumer)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(map[string]interface{}{
									"acctNum":    existingConsumer.AccountNumber,
									"firstName":  existingConsumer.FirstName,
									"middleName": existingConsumer.MiddleName,
									"lastName":   existingConsumer.LastName,
									"suffix":     existingConsumer.Suffix,
									"birthDate":  existingConsumer.BirthDate.Format("2006-01-02"),
									"province":   existingConsumer.Province,
									"postalCode": existingConsumer.PostalCode,
									"cityMun":    existingConsumer.CityMunicipality,
									"barangay":   existingConsumer.Barangay,
									"street":     existingConsumer.Street,
									"phoneNum":   existingConsumer.PhoneNumber,
								})
								return
							}

							// Delete handler
							accNo := r.FormValue("delete-consumer-acc-no")

							// Parse meter number
							accNoInt, err := strconv.Atoi(accNo)
							if err != nil {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Invalid Account number format",
								})
								return
							}

							// Delete document using meter number as ID
							deleteResult, err := svc.DeleteOne(ctx, "consumers", bson.M{"_id": accNoInt})
							if err != nil {
								logger.Sugar().Errorf("Delete failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							if deleteResult.DeletedCount == 0 {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusNotFound)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Consumer Account not found: " + accNo,
								})
								return
							}

							logger.Sugar().Infof("Deleted %d document(s)", deleteResult.DeletedCount)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success": true,
								"message": "Consumer Account deleted successfully",
							})

						default:
							http.NotFound(w, r)
						}

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
					web.SystemAdminEmployeeAccountingWebPage().Render(r.Context(), w)
					// web.SystemAdminEmployeeAccountingWebPage(
					// 	models.AccountingRatesTableFormType.Display,
					// 	models.AccountingRatesTable{
					// 		Date:        "01/01/01",
					// 		Particulars: "RESIDENTIAL",
					// 		Rates:       "",
					// 		ERC:         "9.9298",
					// 		AccountingRatesTableRowGroup: []models.AccountingRatesTableRowGroup{
					// 			// Generation Charges
					// 			{
					// 				Particulars: "Generation Charges",
					// 				Unit:        "",
					// 				Rates:       "5.6092",
					// 				ERC:         "5.6092",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Generation Energy Charge",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "5.6092",
					// 						ERC:         "5.6092",
					// 					},
					// 					{
					// 						Particulars: "Other Generation Rate Adjustment",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.0000",
					// 						ERC:         "0.0000",
					// 					},
					// 				},
					// 			},
					// 			// Transmission Charges
					// 			{
					// 				Particulars: "Transmission Charges (NCCP)",
					// 				Unit:        "",
					// 				Rates:       "0.6853",
					// 				ERC:         "0.6853",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Transmission Demand Charge",
					// 						Unit:        "PhP/kW",
					// 						Rates:       "0.0000",
					// 						ERC:         "0.0000",
					// 					},
					// 					{
					// 						Particulars: "Transmission System Charge",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.6853",
					// 						ERC:         "0.6853",
					// 					},
					// 				},
					// 			},
					// 			// System Loss Charge
					// 			{
					// 				Particulars: "System Loss Charge",
					// 				Unit:        "",
					// 				Rates:       "0.9344",
					// 				ERC:         "0.9344",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "System Loss Charge",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.9344",
					// 						ERC:         "0.9344",
					// 					},
					// 				},
					// 			},
					// 			// Continue with other sections following the same pattern
					// 			// Distribution Charges
					// 			{
					// 				Particulars: "Distribution Charges",
					// 				Unit:        "",
					// 				Rates:       "0.4613",
					// 				ERC:         "0.4613",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Distribution Demand Charge",
					// 						Unit:        "PhP/kW",
					// 						Rates:       "0.0000",
					// 						ERC:         "0.0000",
					// 					},
					// 					{
					// 						Particulars: "Distribution System Charge",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.4613",
					// 						ERC:         "0.4613",
					// 					},
					// 				},
					// 			},
					// 			// Supply Charges
					// 			{
					// 				Particulars: "Supply Charges",
					// 				Unit:        "",
					// 				Rates:       "0.5376",
					// 				ERC:         "0.5376",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Supply Retail Customer Charge",
					// 						Unit:        "PhP/Cust/Mo",
					// 						Rates:       "0.0000",
					// 						ERC:         "0.0000",
					// 					},
					// 					{
					// 						Particulars: "Supply System Charge",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.5376",
					// 						ERC:         "0.5376",
					// 					},
					// 				},
					// 			},
					// 			// Add remaining sections following the same structure...
					// 			// Example for VAT section:
					// 			{
					// 				Particulars: "VAT",
					// 				Unit:        "",
					// 				Rates:       "1.0943",
					// 				ERC:         "0.8543",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Generation",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.6376",
					// 						ERC:         "0.6376",
					// 					},
					// 					{
					// 						Particulars: "Transmission",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.1096",
					// 						ERC:         "0.1096",
					// 					},
					// 					// Add other VAT components...
					// 				},
					// 			},
					// 			// Universal Charge
					// 			{
					// 				Particulars: "Universal Charge",
					// 				Unit:        "",
					// 				Rates:       "0.2250",
					// 				ERC:         "0.2250",
					// 				SubRowGroup: []models.SubRowGroup{
					// 					{
					// 						Particulars: "Missionary Electrification",
					// 						Unit:        "PhP/kWh",
					// 						Rates:       "0.1822",
					// 						ERC:         "0.1822",
					// 					},
					// 					// Add other universal charge components...
					// 				},
					// 			},
					// 		},
					// 	},
					// ).Render(r.Context(), w)
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
					case "rates-data":
						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						// Define the aggregation pipeline
						pipeline := []bson.M{
							{"$match": bson.M{"type": bson.M{"$in": []string{"RATES", "ERC", "INTEREST"}}}},
							{"$group": bson.M{"_id": nil, "docs": bson.M{"$push": "$$ROOT"}}},
							{"$addFields": bson.M{
								"ratesDoc": bson.M{"$arrayElemAt": []interface{}{
									bson.M{"$filter": bson.M{
										"input": "$docs",
										"cond":  bson.M{"$eq": []interface{}{"$$this.type", "RATES"}},
									}},
									0,
								}},
								"ercDoc": bson.M{"$arrayElemAt": []interface{}{
									bson.M{"$filter": bson.M{
										"input": "$docs",
										"cond":  bson.M{"$eq": []interface{}{"$$this.type", "ERC"}},
									}},
									0,
								}},
								"interestDoc": bson.M{"$arrayElemAt": []interface{}{
									bson.M{"$filter": bson.M{
										"input": "$docs",
										"cond":  bson.M{"$eq": []interface{}{"$$this.type", "INTEREST"}},
									}},
									0,
								}},
							}},
							{"$addFields": bson.M{
								"ercTotal": bson.M{"$let": bson.M{
									"vars": bson.M{"totalSection": bson.M{"$arrayElemAt": []interface{}{
										bson.M{"$filter": bson.M{
											"input": "$ercDoc.ratesdata.sections",
											"cond":  bson.M{"$eq": []interface{}{"$$this.name", "TOTAL RATE"}},
										}},
										0,
									}}},
									"in": "$$totalSection.erc",
								}},
								"totalRate": bson.M{"$let": bson.M{
									"vars": bson.M{"totalSection": bson.M{"$arrayElemAt": []interface{}{
										bson.M{"$filter": bson.M{
											"input": "$ratesDoc.ratesdata.sections",
											"cond":  bson.M{"$eq": []interface{}{"$$this.name", "TOTAL RATE"}},
										}},
										0,
									}}},
									"in": "$$totalSection.rate",
								}},
							}},
							{"$project": bson.M{
								"_id":             0,
								"billingDate":     "$ratesDoc.ratesdata.date",
								"type":            "RESIDENTIAL",
								"overdueInterest": "$interestDoc.interestdata.interest",
								"sections": bson.M{"$concatArrays": []interface{}{
									[]interface{}{bson.M{
										"id":   "header-residential",
										"type": "main-header",
										"name": "RESIDENTIAL",
										"rate": "",
										"erc":  "$ercTotal",
									}},
									bson.M{"$map": bson.M{
										"input": bson.M{"$range": []interface{}{0, bson.M{"$size": "$ratesDoc.ratesdata.sections"}, 1}},
										"as":    "idx",
										"in": bson.M{"$let": bson.M{
											"vars": bson.M{
												"rSec": bson.M{"$arrayElemAt": []interface{}{"$ratesDoc.ratesdata.sections", "$$idx"}},
												"eSec": bson.M{"$arrayElemAt": []interface{}{"$ercDoc.ratesdata.sections", "$$idx"}},
											},
											"in": bson.M{"$cond": []interface{}{
												bson.M{"$ne": []interface{}{"$$rSec.name", "TOTAL RATE"}},
												bson.M{"$let": bson.M{
													"vars": bson.M{"sectionId": bson.M{"$switch": bson.M{
														"branches": []bson.M{
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Generation Charges"}}, "then": "cat-gen"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Transmission Charges (NGCP)"}}, "then": "cat-trans"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "System Loss Charge"}}, "then": "cat-sysloss"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Distribution Charges"}}, "then": "cat-dist"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Supply Charges"}}, "then": "cat-supply"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Metering Charges"}}, "then": "cat-meter"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Reinvestment Fund/MCC"}}, "then": "cat-reinvest"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Other Charges"}}, "then": "cat-other"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Franchise Tax"}}, "then": "cat-franchise"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Business Tax"}}, "then": "cat-business"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Real Property Tax"}}, "then": "cat-property"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "VAT"}}, "then": "cat-vat"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "Universal Charge"}}, "then": "cat-universal"},
															{"case": bson.M{"$eq": []interface{}{"$$rSec.name", "FIT - ALL"}}, "then": "cat-fit"},
														},
														"default": "",
													}}},
													"in": bson.M{
														"id":   "$$sectionId",
														"type": "category",
														"name": "$$rSec.name",
														"rate": "$$rSec.rate",
														"erc":  "$$eSec.erc",
														"items": bson.M{"$map": bson.M{
															"input": bson.M{"$range": []interface{}{0, bson.M{"$size": "$$rSec.items"}, 1}},
															"as":    "j",
															"in": bson.M{"$let": bson.M{
																"vars": bson.M{
																	"rItem": bson.M{"$arrayElemAt": []interface{}{"$$rSec.items", "$$j"}},
																	"eItem": bson.M{"$arrayElemAt": []interface{}{"$$eSec.items", "$$j"}},
																},
																"in": bson.M{
																	"id": bson.M{"$concat": []interface{}{
																		"item-",
																		bson.M{"$substrCP": []interface{}{
																			"$$sectionId",
																			4,
																			bson.M{"$subtract": []interface{}{bson.M{"$strLenCP": "$$sectionId"}, 4}},
																		}},
																		bson.M{"$toString": bson.M{"$add": []interface{}{"$$j", 1}}},
																	}},
																	"name": "$$rItem.name",
																	"unit": "$$rItem.unit",
																	"rate": "$$rItem.rate",
																	"erc":  "$$eItem.erc",
																},
															}},
														}},
													},
												}},
												nil,
											}},
										}},
									}},
									[]interface{}{bson.M{
										"id":   "total-section",
										"type": "total",
										"name": "TOTAL RATE",
										"rate": "$totalRate",
										"erc":  "$ercTotal",
									}},
								}},
							}},
							{"$addFields": bson.M{
								"sections": bson.M{"$filter": bson.M{
									"input": "$sections",
									"as":    "sec",
									"cond":  bson.M{"$ne": []interface{}{"$$sec", nil}},
								}},
							}},
						}

						cursor, err := svc.Aggregation(ctx, "rates", pipeline)
						if err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Aggregation failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}
						defer cursor.Close(ctx)

						var result bson.M
						if cursor.Next(ctx) {
							if err := cursor.Decode(&result); err != nil {
								c.Deps.GetLogger().Sugar().Errorf("Decoding failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Failed to decode result",
								})
								return
							}
						} else {
							c.Deps.GetLogger().Sugar().Errorf("No results found")
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusNotFound)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "No data found",
							})
							return
						}

						c.Deps.GetLogger().Sugar().Debugf("Aggregation result: %v", result)

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(result)

					default:
						http.NotFound(w, r)
					}
				case "POST":
					switch formType {
					case "submit-update-rates-form":

						// Decode incoming JSON into “payload”
						var payload models.RatesDocument
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						updateResult, err := svc.UpdateOne(ctx, "rates", bson.M{"type": payload.Type}, bson.M{"$set": payload})
						if err != nil {
							logger.Sugar().Errorf("Update failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}

						if updateResult.MatchedCount == 0 {
							// Insert if no document matched
							insertResult, err := svc.InsertOne(ctx, "rates", payload)
							if err != nil {
								c.Deps.GetLogger().Sugar().Errorf("Insert failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success":     true,
								"type":        "RATES",
								"data":        payload.RatesData,
								"inserted_id": insertResult.InsertedID,
							})

							return
						}

						logger.Sugar().Infof("Updated %d document(s)\tID: %d", updateResult.ModifiedCount, updateResult.UpsertedID)

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"success":     true,
							"type":        "RATES",
							"data":        payload.RatesData,
							"inserted_id": updateResult.UpsertedID,
						})

					case "submit-update-erc-form":

						// Decode incoming JSON into “payload”
						var payload models.RatesDocument
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						updateResult, err := svc.UpdateOne(ctx, "rates", bson.M{"type": payload.Type}, bson.M{"$set": payload})
						if err != nil {
							logger.Sugar().Errorf("Update failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}

						if updateResult.MatchedCount == 0 {
							// Insert if no document matched
							insertResult, err := svc.InsertOne(ctx, "rates", payload)
							if err != nil {
								c.Deps.GetLogger().Sugar().Errorf("Insert failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success":     true,
								"type":        "RATES",
								"data":        payload.RatesData,
								"inserted_id": insertResult.InsertedID,
							})

							return
						}

						logger.Sugar().Infof("Updated %d document(s)\tID: %d", updateResult.ModifiedCount, updateResult.UpsertedID)

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"success":     true,
							"type":        "RATES",
							"data":        payload.RatesData,
							"inserted_id": updateResult.UpsertedID,
						})

					case "submit-update-interest-form":

						// Decode incoming JSON into “payload”
						var payload models.RatesDocument
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						updateResult, err := svc.UpdateOne(ctx, "rates", bson.M{"type": payload.Type}, bson.M{"$set": payload})
						if err != nil {
							logger.Sugar().Errorf("Update failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}

						if updateResult.MatchedCount == 0 {
							// Insert if no document matched
							insertResult, err := svc.InsertOne(ctx, "rates", payload)
							if err != nil {
								c.Deps.GetLogger().Sugar().Errorf("Insert failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Internal server error",
								})
								return
							}

							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"success":     true,
								"type":        "ERC",
								"data":        payload.RatesData,
								"inserted_id": insertResult.InsertedID,
							})

							return
						}

						logger.Sugar().Infof("Updated %d document(s)\tID: %d", updateResult.ModifiedCount, updateResult.UpsertedID)

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"success":     true,
							"type":        "INTEREST",
							"data":        payload.RatesData,
							"inserted_id": updateResult.UpsertedID,
						})

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

				case "POST":
					switch formType {
					case "verify-account-number":
						type accountRequest struct {
							AccountNumber string `json:"accountNumber"`
						}

						var payload accountRequest
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}

						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// Convert accountNumber to int
						accNumInt, err := strconv.Atoi(payload.AccountNumber)
						if err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Invalid account number", http.StatusBadRequest)
							return
						}

						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						pipeline := []bson.M{
							{"$match": bson.M{"acctNum": accNumInt}},
							{
								"$lookup": bson.M{
									"from":         "balances",
									"localField":   "acctNum",
									"foreignField": "acctNum",
									"as":           "balances",
								},
							},
							{"$unwind": bson.M{"path": "$balances", "preserveNullAndEmptyArrays": true}},
							{
								"$project": bson.M{
									"_id": 1,

									"acctNum": 1,
									"consumerName": bson.M{
										"$concat": []string{"$firstName", " ", "$lastName"},
									},
									"consumerType":    "Commercial", // adjust if you have a real source field
									"isActive":        "$isActive",
									"lastPaymentDate": time.Now().UTC(),

									// -------------------------------
									// CURRENT BILL (unchanged)
									// -------------------------------
									"currentBill": bson.M{
										"$cond": bson.M{
											"if": "$balances.currentBill",
											"then": bson.M{
												"billId": "$balances.currentBill.billId",
												"issueDate": bson.M{
													"$dateToString": bson.M{
														"format": "%Y-%m-%dT%H:%M:%SZ",
														"date":   "$balances.currentBill.issueDate",
													},
												},
												"dueDate": bson.M{
													"$dateToString": bson.M{
														"format": "%Y-%m-%dT%H:%M:%SZ",
														"date":   "$balances.currentBill.dueDate",
													},
												},
												"duration": bson.M{
													"start": bson.M{
														"$dateToString": bson.M{
															"format": "%Y-%m-%dT%H:%M:%SZ",
															"date":   "$balances.currentBill.duration.start",
														},
													},
													"end": bson.M{
														"$dateToString": bson.M{
															"format": "%Y-%m-%dT%H:%M:%SZ",
															"date":   "$balances.currentBill.duration.end",
														},
													},
												},
												"charges": bson.M{
													"amountDue":   "$balances.currentBill.charges.amountDue",
													"usedKwH":     "$balances.currentBill.charges.usedKwH",
													"rates":       bson.M{}, // leave empty or reshape if needed
													"overdueFees": nil,
												},
												"isPaid": "$balances.currentBill.isPaid",
											},
											"else": nil,
										},
									},

									// -------------------------------
									// OVERDUE BILL
									// -------------------------------
									"overdueBill": bson.M{
										"$cond": bson.M{
											"if": bson.M{
												// If balances.overdueBill is missing, treat as []
												"$gt": []interface{}{
													bson.M{
														"$size": bson.M{
															"$ifNull": []interface{}{"$balances.overdueBill", []interface{}{}},
														},
													},
													0,
												},
											},
											"then": bson.M{
												"$map": bson.M{
													"input": "$balances.overdueBill",
													"as":    "bill",
													"in": bson.M{
														"billId": "$$bill.billId",
														"issueDate": bson.M{
															"$dateToString": bson.M{
																"format": "%Y-%m-%dT%H:%M:%SZ",
																"date":   "$$bill.issueDate",
															},
														},
														"dueDate": bson.M{
															"$dateToString": bson.M{
																"format": "%Y-%m-%dT%H:%M:%SZ",
																"date":   "$$bill.dueDate",
															},
														},
														"duration": bson.M{
															"start": bson.M{
																"$dateToString": bson.M{
																	"format": "%Y-%m-%dT%H:%M:%SZ",
																	"date":   "$$bill.duration.start",
																},
															},
															"end": bson.M{
																"$dateToString": bson.M{
																	"format": "%Y-%m-%dT%H:%M:%SZ",
																	"date":   "$$bill.duration.end",
																},
															},
														},
														"charges": bson.M{
															"amountDue":   "$$bill.charges.amountDue",
															"usedKwH":     "$$bill.charges.usedKwH",
															"rates":       "$$bill.charges.rates",
															"overdueFees": "$$bill.charges.overdueFees",
														},
														"isPaid": "$$bill.isPaid",
													},
												},
											},
											"else": []interface{}{},
										},
									},

									// -------------------------------
									// BILL HISTORY
									// -------------------------------
									"billHistory": bson.M{
										"$cond": bson.M{
											"if": bson.M{
												"$gt": []interface{}{
													bson.M{
														"$size": bson.M{
															"$ifNull": []interface{}{"$balances.billHistory", []interface{}{}},
														},
													},
													0,
												},
											},
											"then": bson.M{
												"$map": bson.M{
													"input": "$balances.billHistory",
													"as":    "bill",
													"in": bson.M{
														"billId": "$$bill.billId",
														"issueDate": bson.M{
															"$dateToString": bson.M{
																"format": "%Y-%m-%dT%H:%M:%SZ",
																"date":   "$$bill.issueDate",
															},
														},
														"dueDate": bson.M{
															"$dateToString": bson.M{
																"format": "%Y-%m-%dT%H:%M:%SZ",
																"date":   "$$bill.dueDate",
															},
														},
														"duration": bson.M{
															"start": bson.M{
																"$dateToString": bson.M{
																	"format": "%Y-%m-%dT%H:%M:%SZ",
																	"date":   "$$bill.duration.start",
																},
															},
															"end": bson.M{
																"$dateToString": bson.M{
																	"format": "%Y-%m-%dT%H:%M:%SZ",
																	"date":   "$$bill.duration.end",
																},
															},
														},
														"charges": bson.M{
															"amountDue":   "$$bill.charges.amountDue",
															"usedKwH":     "$$bill.charges.usedKwH",
															"rates":       "$$bill.charges.rates",
															"overdueFees": "$$bill.charges.overdueFees",
														},
														"isPaid": "$$bill.isPaid",
													},
												},
											},
											"else": []interface{}{},
										},
									},

									// -------------------------------
									// PAYMENT HISTORY
									// -------------------------------
									"paymentHistory": bson.M{
										"$cond": bson.M{
											"if": bson.M{
												"$gt": []interface{}{
													bson.M{
														"$size": bson.M{
															"$ifNull": []interface{}{"$balances.paymentHistory", []interface{}{}},
														},
													},
													0,
												},
											},
											"then": bson.M{
												"$map": bson.M{
													"input": "$balances.paymentHistory",
													"as":    "payment",
													"in": bson.M{
														"transactionId": "$$payment.transactionId",
														"billIds":       "$$payment.billIds",
														"date": bson.M{
															"$dateToString": bson.M{
																"format": "%Y-%m-%dT%H:%M:%SZ",
																"date":   "$$payment.date",
															},
														},
														"amount": "$$payment.amount",
													},
												},
											},
											"else": []interface{}{},
										},
									},
								},
							},
						}

						cursor, err := svc.Aggregation(ctx, "consumers", pipeline)
						if err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Aggregation failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}
						defer cursor.Close(ctx)

						var result bson.M
						if cursor.Next(ctx) {
							if err := cursor.Decode(&result); err != nil {
								c.Deps.GetLogger().Sugar().Errorf("Decoding failed: %v", err)
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]string{
									"error": "Failed to decode result",
								})
								return
							}
						} else {
							c.Deps.GetLogger().Sugar().Errorf("No results found")
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusNotFound)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "No data found",
							})
							return
						}

						c.Deps.GetLogger().Sugar().Debugf("Aggregation result: %v", result)

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(result)

					case "process-payment":
						type ProcessPayment struct {
							AccountNumber string   `json:"accountNumber"`
							TotalAmount   float64  `json:"amount"`
							BillIds       []string `json:"billIds"`
						}

						var payload ProcessPayment
						if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
							c.Deps.GetLogger().Sugar().Errorf("Decode error: %v", err)
							http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
							return
						}
						c.Deps.GetLogger().Sugar().Infof("Received valid payload: %+v", payload)

						// Process Data
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						svc := database.New()

						payment := models.PaymentHistory{
							TransactionId: internal.GenerateUUIDBillingID(),
							BillIds:       payload.BillIds,
							Date:          time.Now().UTC(),
							Amount:        payload.TotalAmount,
						}

						updateData := bson.M{
							"$set": bson.M{
								"currentBill.isPaid": true,
							},
							"$push": bson.M{
								"paymentHistory": bson.M{
									"$each": []bson.M{
										{
											"transactionId": payment.TransactionId,
											"billIds":       payment.BillIds,
											"date":          payment.Date,
											"amount":        payment.Amount,
										},
									},
								},
							},
						}

						accNoInt, err := strconv.Atoi(payload.AccountNumber)
						if err != nil {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Invalid Account number format",
							})
							return
						}

						updateResult, err := svc.UpdateOne(ctx, "balances", bson.M{"acctNum": accNoInt}, updateData)
						if err != nil {
							logger.Sugar().Errorf("Update failed: %v", err)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{
								"error": "Internal server error",
							})
							return
						}

						if updateResult.MatchedCount == 0 {
							// return an error
							w.WriteHeader(http.StatusInternalServerError)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"errro": "No Matched Document",
							})
							return
						}

						logger.Sugar().Infof("Updated %d document(s)", updateResult.ModifiedCount)

						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"amount":        payment.Amount,
							"accountNumber": payload.AccountNumber,
							"date":          payment.Date,
							"transactionId": payment.TransactionId,
						})

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
