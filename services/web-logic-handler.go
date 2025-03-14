/* package services

import (
	"fmt"
	"log"
	"net/http"
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

	store := utils.GetStore(config.MaxWebUsers, config.MaxWebSessions, config.MaxMeterUsers, config.MaxMeterSessions)
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
	store := utils.GetStore(config.MaxWebUsers, config.MaxWebSessions, config.MaxMeterUsers, config.MaxMeterSessions)

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
