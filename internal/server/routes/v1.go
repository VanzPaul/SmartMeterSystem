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
	"fmt"
	"log"
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

							// Validate required fields
							if form_meterNo == "" || form_meterInstallationDate == "" || form_consumerTransformerId == "" ||
								form_meterLatitude == "" || form_meterLongitude == "" || form_consumerAccNo == "" {
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

							bson, err_bson := bson.Marshal(models.MeterDocument{
								ID:                    meterNo,
								ConsumerAccNo:         consumerAccNo,
								InstallationDate:      meterInstallationDate,
								ConsumerTransformerID: consumerTransformerId,
								Latitude:              meterLatitude,
								Longitude:             meterLongitude,
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
							insertResult, err := svc.InsertOne(ctx, "meters", bson)

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

							bson, err_bson := bson.Marshal(models.ConsumerDocument{
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
							insertResult, err := svc.InsertOne(ctx, "consumers", bson)

							if err != nil {
								// Check for duplicate key error
								var writeException mongo.WriteException
								if errors.As(err, &writeException) {
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

						default:
							http.NotFound(w, r)
						}

					case "update":
						if len(pathParts) < 2 {
							http.NotFound(w, r)
							return
						}

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
									"latitude":         existingMeter.Latitude,
									"longitude":        existingMeter.Longitude,
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
									"latitude":         existingMeter.Latitude,
									"longitude":        existingMeter.Longitude,
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
