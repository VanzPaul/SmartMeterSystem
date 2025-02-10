package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/utils"
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

	// FIXME: this should be a database find call
	store := models.GetStore() // Use the singleton store
	if _, ok := store.Users[email]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		config.Logger.Fatal("password hashing failed: ",
			zap.Error(err))
	}

	// Make db call to find the role of the account number
	// accountNoData := ""

	// Create account struct instance
	// acc := models.Account{
	// 	HashedPassword: hashedPassword,
	// 	Email: email,
	// 	CreatedAt: time.Now().UTC().Unix(),
	// 	UpdatedAt: time.Now().UTC().Unix(),
	// 	Role: role,
	// 	Status: models.AccountStatus{
	// 		IsActive: true,
	// 	},
	// 	RoleSpecificDataID: ,

	// }
	// Get the context and pass this to the database functions
	//ctx := r.Context()
	// Use context in database operations
	// err := createDocument(ctx, document)
	// if err != nil {
	// 	handleError(w, r, err)
	// 	return
	// }

	// TODO: create a document in mongodb using services/createDocument()
	// remove this
	store.Users[email] = models.LoginData{
		HashedPassword: hashedPassword,
		AccountNo:      accountNo,
	}

	// remove this
	// Convert the user struct to JSON
	jsonData, _ := json.MarshalIndent(store.Users[email], "", "  ")
	log.Printf("store.Users[%s]: %s", email, string(jsonData))

	fmt.Fprintln(w, "User Registered Successfully")
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
