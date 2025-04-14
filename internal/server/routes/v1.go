/*
 * @file internal/server/routes/v1.go
 * @brief v1.go file holds all v1 route groups and their handlers
 */
package routes

import (
	"SmartMeterSystem/cmd/web"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.NotFound().Render(r.Context(), w)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)

		// Inside your handler function
		userType := r.URL.Query().Get("user_type")

		web.LoginWebPage(c.Deps.GetDefaultRouteVersion(), userType).Render(r.Context(), w)
	})

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.ConsumerDashboardWebPage().Render(r.Context(), w)
	})

	return mux
}

func (c *V1EmployeeRoute) HandleV1() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.NotFound().Render(r.Context(), w)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)

		// Inside your handler function
		userType := r.URL.Query().Get("user_type")

		web.LoginWebPage(c.Deps.GetDefaultRouteVersion(), userType).Render(r.Context(), w)
	})

	mux.HandleFunc("/login/submit", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		w.Header().Set("HX-Redirect", "/v1/employee/sysadmin/dashboard") // REMINDER: Make the redirect Dynamic
		w.WriteHeader(http.StatusOK)
	})

	// System Admin Logout Route
	mux.HandleFunc("/sysadmin/logout", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		w.Header().Set("HX-Redirect", "/home")
		w.WriteHeader(http.StatusOK)
	})

	// System Admin Dashboard Routes
	mux.HandleFunc("/sysadmin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.SystemAdminEmployeeDashboardWebPage().Render(r.Context(), w)
	})

	// System Admin Account Routes
	mux.HandleFunc("/sysadmin/accounts", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.SystemAdminEmployeeAccountsWebPage().Render(r.Context(), w)
	})
	mux.HandleFunc("/sysadmin/accounts/meter-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.NewMeterAccountForm().Render(r.Context(), w)
	})
	mux.HandleFunc("/sysadmin/accounts/consumer-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.NewConsumerAccountForm().Render(r.Context(), w)
	})
	mux.HandleFunc("/sysadmin/accounts/employee-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.NewEmployeeAccountForm().Render(r.Context(), w)
	})

	// System Admin Consumer Routes
	mux.HandleFunc("/sysadmin/consumer", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		web.SystemAdminEmployeeConsumerWebPage().Render(r.Context(), w)
	})

	// System Admin Accounting Routes
	mux.HandleFunc("/sysadmin/accounting", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
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
	})

	// System Admin Accounting Routes
	mux.HandleFunc("/sysadmin/accounting/update-rates-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		// web.SystemAdminEmployeeAccountingWebPage().Render(r.Context(), w)
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
	})

	mux.HandleFunc("/sysadmin/accounting/submit-update-rates-form", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %v", r), http.StatusInternalServerError)
				log.Printf("Panic recovered: %v", r)
			}
		}()

		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)

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
	})

	mux.HandleFunc("/sysadmin/accounting/update-erc-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)
		// web.SystemAdminEmployeeAccountingWebPage().Render(r.Context(), w)
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
	})

	mux.HandleFunc("/sysadmin/accounting/submit-update-erc-form", func(w http.ResponseWriter, r *http.Request) {
		c.Deps.GetLogger().Sugar().Infof("%s\t%s\t%s", c.Deps.GetDefaultRouteVersion(), r.URL, r.Method)

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
	})

	return mux
}
