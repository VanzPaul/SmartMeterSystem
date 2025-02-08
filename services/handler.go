package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	accountNo := r.FormValue("accountNo")

	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	store := models.GetStore() // Use the singleton store
	if _, ok := store.Users[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	store.Users[username] = models.LoginData{
		HashedPassword: hashedPassword,
		AccountNo:      accountNo,
	}

	fmt.Fprintln(w, "User Registered Successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	store := models.GetStore() // Use the singleton store
	user, ok := store.Users[username]
	if !ok || !utils.CheckPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid Username or Password", http.StatusUnauthorized)
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
	store.Users[username] = user
	store.Sessions[sessionToken] = username // Map session token to username

	log.Printf("Session Token: %s\nCSRF Token: %s\n", user.SessionToken, user.CSRFToken)
	fmt.Fprintln(w, "Login Successful!")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s\n", username)
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
		username := store.Sessions[sessionToken.Value]
		delete(store.Sessions, sessionToken.Value) // Remove session token
		user := store.Users[username]
		user.SessionToken = ""
		user.CSRFToken = ""
		store.Users[username] = user
	}

	fmt.Fprintln(w, "Logged Out Successfully!")
}
