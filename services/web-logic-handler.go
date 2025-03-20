/* package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// Register creates a new account
func SubmitRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	accountNo := r.FormValue("accountNo")
	if len(email) < 4 || len(password) < 8 {
		http.Error(w, "Invalid email/password", http.StatusNotAcceptable)
		return
	}

	ctx := r.Context()
	utils.Logger.Debug("Creating new MongoDB Controller")
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Check if email already exists
	var emailData []models.Account
	emailFilter := bson.M{"email": email}
	emailfindErr := db.Find(ctx, models.Accounts, emailFilter, &emailData)
	if emailfindErr != nil {
		utils.Logger.Error("Checking if email exist error", zap.Any("Error", emailfindErr))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(emailData) != 0 {
		fmt.Fprintf(w, `<p class="error">Email does not exist.</p>`)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.Logger.Fatal("password hashing failed: ", zap.Error(err))
	}

	// Rest of your registration logic
	var consumer client.Consumer
	filter := bson.M{"accountNumber": accountNo}
	findErr := db.FindOne(ctx, models.Consumers, filter, &consumer)
	if findErr != nil {
		utils.Logger.Sugar().Errorf("Error finding account: %v", findErr)
		http.Error(w, "Invalid account number", http.StatusBadRequest)
		return
	}

	account := models.Account{
		HashedPassword: hashedPassword,
		Email:          email,
		CreatedAt:      time.Now().UTC().Unix(),
		UpdatedAt:      time.Now().UTC().Unix(),
		Role:           models.Role(consumer.AccountType),
		Status: models.AccountStatus{
			IsActive: true,
		},
		RoleSpecificDataID: consumer.ID,
	}

	insertResult, createErr := CreateDocument(ctx, db, models.Accounts, &account)
	if createErr != nil {
		log.Fatalf("Err creating document %s: %v\n", models.Accounts, createErr)
	}
	utils.Logger.Sugar().Debugf("Insert Successful: %s", insertResult.String())

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

// Login authenticates a user
func SubmitLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	utils.Logger.Sugar().Debugf("email: %s\tpassword: %s\n", email, password)

	store := utils.GetStore(utils.MaxWebUsers, utils.MaxWebSessions, utils.MaxMeterUsers, utils.MaxMeterSessions)
	ctx := r.Context()
	utils.Logger.Debug("Creating new MongoDB Controller")

	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to create MongoDB controller: %v", err)
		fmt.Fprintf(w, `<p class="error">Internal server error. Please try again later.</p>`)
		return
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			utils.Logger.Sugar().Errorf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Fetch and check account hashed password
	var account models.Account
	filter := bson.M{"email": email}
	accErrs := db.FindOne(ctx, models.Accounts, filter, &account)
	if accErrs != nil {
		utils.Logger.Sugar().Errorf("Error finding account: %v", accErrs)
		fmt.Fprintf(w, `<p class="error">Email does not exist.</p>`)
		return
	}

	if !utils.CheckPasswordHash(password, string(account.HashedPassword)) {
		utils.Logger.Sugar().Errorf("Invalid password for email: %s", email)
		fmt.Fprintf(w, `<p class="error">Invalid password. Please try again.</p>`)
		return
	}

	utils.Logger.Debug("Email and Password hash matched!")

	// Generate tokens
	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
	})

	// Store session data
	store.AddWebUser(email, utils.LoginData{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		Role:         string(account.Role),
	})
	store.AddWebSession(sessionToken, email)

	w.Header().Set("HX-Redirect", "/client/dashboard")
	w.WriteHeader(http.StatusOK)
}

// Logout clears the session and CSRF tokens
func Logout(w http.ResponseWriter, r *http.Request) {
	store := models.GetStore(utils.MaxWebUsers, utils.MaxWebSessions, utils.MaxMeterUsers, utils.MaxMeterSessions)

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// Clear tokens from store
	sessionToken, err := r.Cookie("session_token")
	if err == nil {
		if username, exists := store.GetWebSession(sessionToken.Value); exists {
			store.DeleteWebSession(sessionToken.Value)
			if userData, exists := store.GetWebUser(username); exists {
				userData.SessionToken = ""
				userData.CSRFToken = ""
				store.AddWebUser(username, userData)
			}
		}
	}

	fmt.Fprintln(w, "Logged Out Successfully!")
}
*/

package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/client"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

/*************************************************************************************************************/
/*----------------------------------------------  UTILITIES  ------------------------------------------------*/
/*************************************************************************************************************/

// <<<<-------------------------------------------------------->>>> //
// <<<<-------------- Login authenticates a user -------------->>>> //
// <<<<-------------------------------------------------------->>>> //

// SubmitWebLogin handles user authentication for web login.
//
// This function processes the login request by validating the provided email and password,
// checking the credentials against the database, generating session and CSRF tokens,
// and setting cookies for the authenticated user. It also logs relevant debug and error
// information throughout the process.
//
// @param w http.ResponseWriter - The response writer used to send HTTP responses back to the client.
//
//	This includes error messages, redirection headers, and cookies.
//
// @param r *http.Request - The HTTP request object containing the user's login credentials
//
//	(email and password) submitted via form data. It also provides the context for
//	database operations and session management.
func SubmitWebLogin(w http.ResponseWriter, r *http.Request) {
	// Extract email and password from the form data in the request.
	// These are the credentials submitted by the user during login.
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Log the email and password for debugging purposes.
	// Note: Ensure sensitive data like passwords is masked or omitted in production logs.
	utils.Logger.Sugar().Debugf("email: %s\tpassword: %s\n", email, password)

	// Initialize the session store with predefined limits for web users and sessions.
	// This ensures that the application can manage active users and sessions efficiently.
	store := utils.GetStore(config.MaxWebUsers, config.MaxWebSessions, config.MaxMeterUsers, config.MaxMeterSessions)

	// Retrieve the context from the request for database operations.
	// The context allows cancellation and timeout management for long-running operations.
	ctx := r.Context()

	// Log the creation of a new MongoDB controller for database interaction.
	utils.Logger.Debug("Creating new MongoDB Controller")

	// Initialize a MongoDB controller using the provided configuration.
	// This controller is responsible for querying and managing user data in the database.
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		// Log the error if the MongoDB controller fails to initialize.
		// Return an internal server error message to the client.
		utils.Logger.Sugar().Errorf("Failed to create MongoDB controller: %v", err)
		fmt.Fprintf(w, `<p class="error">Internal server error. Please try again later.</p>`)
		return
	}

	// Ensure the MongoDB connection is closed after the function completes.
	// This prevents resource leaks and ensures proper cleanup.
	defer func() {
		if err := db.Close(ctx); err != nil {
			// Log any errors encountered while closing the MongoDB connection.
			utils.Logger.Sugar().Errorf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Define a filter to query the database for the account associated with the provided email.
	var account models.Account
	filter := bson.M{"email": email}

	// Query the database to find the account matching the email.
	// If no account is found, return an error message indicating the email does not exist.
	accErrs := db.FindOne(ctx, models.Accounts, filter, &account)
	if accErrs != nil {
		utils.Logger.Sugar().Errorf("Error finding account: %v", accErrs)
		fmt.Fprintf(w, `<p class="error">Email does not exist.</p>`)
		return
	}

	// Validate the provided password against the hashed password stored in the database.
	// If the password is invalid, return an error message prompting the user to try again.
	if !utils.CheckPasswordHash(password, string(account.HashedPassword)) {
		utils.Logger.Sugar().Errorf("Invalid password for email: %s", email)
		fmt.Fprintf(w, `<p class="error">Invalid password. Please try again.</p>`)
		return
	}

	// Log a success message if the email and password hash match.
	utils.Logger.Debug("Email and Password hash matched!")

	// Generate a session token and a CSRF token for the authenticated user.
	// These tokens are used for session management and CSRF protection, respectively.
	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	// Set a session token cookie for the user.
	// The session token is stored in an HttpOnly cookie to prevent XSS attacks.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Set a CSRF token cookie for the user.
	// The CSRF token is stored in a non-HttpOnly cookie to allow client-side access for validation.
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
	})

	// Store the session data in the in-memory store.
	// This includes the session token, CSRF token, and user role for tracking active sessions.
	store.AddWebUser(email, utils.LoginData{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		Role:         string(account.Role),
	})

	// Associate the session token with the user's email in the session store.
	store.AddWebSession(sessionToken, email)

	// Redirect the user to the dashboard upon successful authentication.
	// The HX-Redirect header is used for redirection in HTMX-based applications.
	w.Header().Set("HX-Redirect", "/client/consumer/dashboard")
	w.WriteHeader(http.StatusOK)
}

// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //

// <<<<-------------------------------------------------------->>>> //
// <<<<------ Logout clears the session and CSRF tokens ------->>>> //
// <<<<-------------------------------------------------------->>>> //

// Logout clears the session and CSRF tokens for a logged-out user.
//
// This function handles the logout process by clearing the session and CSRF cookies,
// removing the associated session data from the in-memory store, and returning a success message.
// It ensures that the user's session is properly terminated and prevents unauthorized access.
//
// @param w http.ResponseWriter - The response writer used to send HTTP responses back to the client.
//
//	This includes setting cookies and returning a logout success message.
//
// @param r *http.Request - The HTTP request object containing the user's session cookie.
//
//	This is used to identify the session to be cleared from the in-memory store.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Initialize the session store with predefined limits for web users and sessions.
	// This ensures that the application can manage active users and sessions efficiently.
	store := utils.GetStore(config.MaxWebUsers, config.MaxWebSessions, config.MaxMeterUsers, config.MaxMeterSessions)

	// Clear the session token cookie by setting its value to an empty string and expiring it immediately.
	// The HttpOnly flag ensures that the cookie cannot be accessed via JavaScript, mitigating XSS risks.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to a past time to invalidate the cookie.
		HttpOnly: true,
	})

	// Clear the CSRF token cookie by setting its value to an empty string and expiring it immediately.
	// The HttpOnly flag is set to false to allow client-side access for validation purposes.
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to a past time to invalidate the cookie.
		HttpOnly: false,
	})

	// Retrieve the session token from the request cookies.
	// This token is used to identify the session to be cleared from the in-memory store.
	sessionToken, err := r.Cookie("session_token")
	if err == nil {
		// Check if the session token exists in the session store and retrieve the associated username.
		if username, exists := store.GetWebSession(sessionToken.Value); exists {
			// Delete the session token from the session store to terminate the session.
			store.DeleteWebSession(sessionToken.Value)

			// Check if the user data exists in the user store and clear the session and CSRF tokens.
			if userData, exists := store.GetWebUser(username); exists {
				userData.SessionToken = ""           // Clear the session token from the user data.
				userData.CSRFToken = ""              // Clear the CSRF token from the user data.
				store.AddWebUser(username, userData) // Update the user data in the store.
			}
		}
	}

	// Return a success message to indicate that the user has been logged out successfully.
	fmt.Fprintln(w, "Logged Out Successfully!")
}

// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //

/*************************************************************************************************************/
/*************************************************************************************************************/
//
//
//
//
//
/*************************************************************************************************************/
/*--------------------------------  API LOGIC FOR CREATING USER ACCOUNTS  -----------------------------------*/
/*************************************************************************************************************/

// <<<<-------------------------------------------------------->>>> //
// <<<<----- Register creates a new general account ----------->>>> //
// <<<<-------------------------------------------------------->>>> //

// CreateGeneralAccount handles the creation of a new user account.
//
// This function processes the registration request by validating the provided email, password,
// and account number, checking for duplicate emails, hashing the password, associating the account
// with a consumer profile, and saving the new account in the database. It also logs relevant debug
// and error information throughout the process.
//
// @param w http.ResponseWriter - The response writer used to send HTTP responses back to the client.
//
//	This includes error messages, redirection headers, and success status codes.
//
// @param r *http.Request - The HTTP request object containing the user's registration details
//
//	(email, password, and account number) submitted via form data. It also provides the context
//	for database operations.
func CreateGeneralAccount(w http.ResponseWriter, r *http.Request) {
	// Extract email, password, and account number from the form data in the request.
	// These are the credentials and details submitted by the user during registration.
	email := r.FormValue("email")
	password := r.FormValue("password")
	accountNo := r.FormValue("accountNo")

	// Perform basic validation on the form data to ensure compliance with requirements.
	// - Email must be at least 4 characters long.
	// - Password must be at least 8 characters long.
	// - Account number cannot be empty.
	if len(email) < 4 {
		http.Error(w, "Invalid email", http.StatusNotAcceptable)
		return
	} else if len(password) < 8 {
		http.Error(w, "Password needs to be at least 8 characters long", http.StatusNotAcceptable)
		return
	} else if len(accountNo) == 0 {
		http.Error(w, "Account number cannot be empty", http.StatusNotAcceptable)
		return
	}

	// Retrieve the context from the request for database operations.
	// The context allows cancellation and timeout management for long-running operations.
	ctx := r.Context()

	// Log the creation of a new MongoDB controller for database interaction.
	utils.Logger.Debug("Creating new MongoDB Controller")

	// Initialize a MongoDB controller using the provided configuration.
	// This controller is responsible for querying and managing user data in the database.
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		// Log a fatal error if the MongoDB controller fails to initialize.
		// This ensures that critical failures are immediately visible.
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}

	// Ensure the MongoDB connection is closed after the function completes.
	// This prevents resource leaks and ensures proper cleanup.
	defer func() {
		if err := db.Close(ctx); err != nil {
			// Log any errors encountered while closing the MongoDB connection.
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Check if the email already exists in the database to prevent duplicate accounts.
	var emailData []models.Account
	emailFilter := bson.M{"email": email}
	emailfindErr := db.Find(ctx, models.Accounts, emailFilter, &emailData)
	if emailfindErr != nil {
		// Log the error if the query to check for an existing email fails.
		// Return an internal server error message to the client.
		utils.Logger.Error("Checking if email exists error", zap.Any("Error", emailfindErr))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(emailData) != 0 {
		// If the email already exists, return an error message indicating the conflict.
		fmt.Fprintf(w, `<p class="error">Email already exists.</p>`)
		return
	}

	// Hash the password securely before storing it in the database.
	// This ensures that sensitive information is protected from unauthorized access.
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		// Log a fatal error if password hashing fails.
		// This ensures that critical security failures are immediately visible.
		utils.Logger.Fatal("password hashing failed: ", zap.Error(err))
	}

	// Query the database to find the consumer profile associated with the provided account number.
	var consumer client.Consumer
	filter := bson.M{"accountNumber": accountNo}
	findErr := db.FindOne(ctx, models.Consumers, filter, &consumer)
	if findErr != nil {
		// Log the error if no consumer profile is found for the given account number.
		// Return an error message indicating the invalid account number.
		utils.Logger.Sugar().Errorf("Error finding account: %v", findErr)
		http.Error(w, "Invalid account number", http.StatusBadRequest)
		return
	}

	// Create a new account object with the validated and processed data.
	account := models.Account{
		HashedPassword: hashedPassword,
		Email:          email,
		CreatedAt:      time.Now().UTC().Unix(),
		UpdatedAt:      time.Now().UTC().Unix(),
		Role:           models.Role(consumer.AccountType),
		Status: models.AccountStatus{
			IsActive: true,
		},
		RoleSpecificDataID: consumer.ID,
	}

	// Insert the new account into the database.
	insertResult, createErr := CreateDocument(ctx, db, models.Accounts, &account)
	if createErr != nil {
		// Log a fatal error if the account creation fails.
		// This ensures that critical failures are immediately visible.
		log.Fatalf("Err creating document %s: %v\n", models.Accounts, createErr)
	}

	// Log the successful insertion of the account into the database.
	utils.Logger.Sugar().Debugf("Insert Successful: %s", insertResult.String())

	// Redirect the user to the login page upon successful registration.
	// The HX-Redirect header is used for redirection in HTMX-based applications.
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //

// <<<<-------------------------------------------------------->>>> //
// <<<<------- Register creates a new role accounts ----------->>>> //
// <<<<-------------------------------------------------------->>>> //

func CreateSystemAdminAccout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to create a new system admin account
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 {
		utils.Logger.Debug("Username cannot be empty")
		http.Error(w, "Username cannot be empty", http.StatusNotAcceptable)
		return
	} else if len(password) == 0 {
		utils.Logger.Debug("Password cannot be empty")
		http.Error(w, "Password cannot be emoty", http.StatusNotAcceptable)
		return
	}

	ctx := r.Context()

	utils.Logger.Debug("Creating new MongoDB Controller")

	// Initialize a MongoDB controller using the provided configuration.
	// This controller is responsible for querying and managing user data in the database.
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		// Log a fatal error if the MongoDB controller fails to initialize.
		// This ensures that critical failures are immediately visible.
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}

	// Ensure the MongoDB connection is closed after the function completes.
	// This prevents resource leaks and ensures proper cleanup.
	defer func() {
		if err := db.Close(ctx); err != nil {
			// Log any errors encountered while closing the MongoDB connection.
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

}

func CreateConsumerAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to create a new consumer account

}

// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //
// <<<<-------------------------------------------------------->>>> //

// TODO: finish this and add authenication
func CreateMeterAccount(w http.ResponseWriter, r *http.Request) {
	sn := r.FormValue("sn")
	snid := r.FormValue("snid")
	transformerid := r.FormValue("transformerid")
	manufacturer := r.FormValue("manufacturer")
	model := r.FormValue("model")
	phase := r.FormValue("phase")
	iccid := r.FormValue("iccid")
	mobilenumber := r.FormValue("mobilenumber")
	longitude := r.FormValue("longitude")
	latitude := r.FormValue("latitude")

	// Convert string coordinates to float64
	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		http.Error(w, "Invalid longitude value", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		http.Error(w, "Invalid latitude value", http.StatusBadRequest)
		return
	}

	// REMINDER: perform generic checking here in the future
	if len(sn) == 0 {
		http.Error(w, "Serial number cannot be empty", http.StatusNotAcceptable)
		return
	} else if len(manufacturer) == 0 {
		http.Error(w, "Manufacturer cannot be empty", http.StatusNotAcceptable)
		return
	} else if len(model) == 0 {
		http.Error(w, "Model cannot be empty", http.StatusNotAcceptable)
		return
	} else if len(phase) == 0 {
		http.Error(w, "Phase cannot be empty", http.StatusNotAcceptable)
	} else if len(snid) == 0 {
		http.Error(w, "SNID cannot be empty", http.StatusNotAcceptable)
	} else if len(transformerid) == 0 {
		http.Error(w, "Transformer ID cannot be empty", http.StatusNotAcceptable)
	}

	ctx := r.Context()
	utils.Logger.Debug("Creating new MongoDB Controller")
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Check if sn already exist
	// REMINDER: temporarily commented for testing purposes
	/* 	var snData []client.Meter
	   	snFilter := bson.M{"sn": sn}
	   	snfindErr := db.Find(ctx, models.Meters, snFilter, &snData)
	   	if snfindErr != nil {
	   		utils.Logger.Error("Checking if meter sn exist error", zap.Any("Error", snfindErr))
	   		http.Error(w, "Internal server error", http.StatusInternalServerError)
	   		return
	   	}

	   	// Return error if existing
	   	if len(snData) == 0 {
	   		fmt.Fprintf(w, `<p class="error">Meter SN already exist.</p>`)
	   		return
	   	} */

	// Generate API-key using the sn
	apikey, err := utils.GenerateAPIKey([]byte(config.MeterEnv.SecretKey), sn)
	if err != nil {
		utils.Logger.Error("Error Generating API key", zap.Any("Error", err))
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Generate Wifi SSID
	wifissid, err := utils.GenerateWifiSSID(sn)
	if err != nil {
		utils.Logger.Error("Error Generating Wifi SSID", zap.Any("Error", err))
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Generate Wifi Password
	wifipassword, err := utils.GenerateWifiPassword(sn)
	if err != nil {
		utils.Logger.Error("Error Generating Wifi Password", zap.Any("Error", err))
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Create meter struct instance
	meter := client.Meter{
		SNID: snid,
		MeterConfig: client.MeterConfig{
			SN:           sn,
			Manufacturer: manufacturer,
			Model:        model,
			Phase:        phase,
			APIKey:       apikey, // FIXME: create a function that generates API key using the SN
			Wifi: client.Wifi{
				SSID:     wifissid,
				Password: wifipassword,
			},
			SIM: client.SIM{
				ICCID:          iccid,
				MobileNumber:   mobilenumber,
				DataUsageMb:    0.0,
				ActivationDate: time.Now().Unix(),
			},
		},
		TransformerID: transformerid,
		Location: client.GeoJSON{
			Type:        "Point",             // Fixed: Changed "point" to "Point"
			Coordinates: []float64{lng, lat}, // Replace with actual coordinates
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
			History: []client.AlertEvent{},
		},
		Commands: client.Commands{
			Active:  []client.ActiveCommand{},
			History: []client.HistoryCommand{},
		},
		Status: client.MeterStatus{
			LastSeen:       0,
			GridConnection: false,
			BatteryLevel:   0.0,
		},
	}

	utils.Logger.Sugar().Debug("New meter:", meter)
	// REMINDER: temporarily commented for testing purposes
	// Insert the new meter document into the database
	/* insertResult, createErr := CreateDocument(ctx, db, models.Meters, &meter)
	if createErr != nil {
		// Log a fatal error if the account creation fails.
		// This ensures that critical failures are immediately visible.
		log.Fatalf("Err creating document %s: %v\n", models.Meters, createErr)
	} */

	// Log the successful insertion of the meter documetn into the database.
	/* utils.Logger.Sugar().Debugf("Insert Successful: %s", insertResult.String()) */

	// Redirect the user to the login page upon successful registration.
	// The HX-Redirect header is used for redirection in HTMX-based applications.
	// REMINDER: uncomment the redirect when the webpage is created
	// w.Header().Set("HX-Redirect", "")
	w.WriteHeader(http.StatusOK)

}

/*************************************************************************************************************/
/*************************************************************************************************************/
