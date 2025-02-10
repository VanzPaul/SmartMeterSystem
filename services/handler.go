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

func Register(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "User already exists", http.StatusConflict)
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
	fmt.Fprintln(w, insertResult.String())
}

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	store := models.GetStore() // Use the singleton store
	user, ok := store.Users[email]
	if !ok || !utils.CheckPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid email or Password", http.StatusUnauthorized)
		return
	}

	// Set session cookie
	sessionToken := utils.GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	// Set CSRF cookie
	csrfToken := utils.GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	// Store session & CSRF token in the database
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	store.Users[email] = user
	store.Sessions[sessionToken] = email // Map session token to email

	log.Printf("Session Token: %s\nCSRF Token: %s\n", user.SessionToken, user.CSRFToken)
	fmt.Fprintln(w, "Login Successful!")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s\n", email)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	store := models.GetStore() // Use the singleton store

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

	// Clear tokens from the database
	sessionToken, err := r.Cookie("session_token")
	if err == nil {
		email := store.Sessions[sessionToken.Value]
		delete(store.Sessions, sessionToken.Value) // Remove session token
		user := store.Users[email]
		user.SessionToken = ""
		user.CSRFToken = ""
		store.Users[email] = user
	}

	fmt.Fprintln(w, "Logged Out Successfully!")
}
