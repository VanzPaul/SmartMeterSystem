package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vanspaul/SmartMeterSystem/utils"
)

type login struct {
	HashedPassword, AccountNo, SessionToken, CSRFToken string
}

// Database. Change this with MongoDB
var users = map[string]login{}
var sessions = map[string]string{} // Map session tokens to usernames

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	accountNo := r.FormValue("accountNo")
	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}
	if _, ok := users[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hashedPassword, _ := utils.HashPassword(password)
	users[username] = login{
		HashedPassword: hashedPassword,
		AccountNo:      accountNo,
	}
	fmt.Fprintln(w, "User Registered Successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, ok := users[username]
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
	users[username] = user
	sessions[sessionToken] = username // Map session token to username
	log.Printf("Session Token: %s\nCSRF Token: %s\n", user.SessionToken, user.CSRFToken)
	fmt.Fprintln(w, "Login Successful!")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s", username)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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
		username := sessions[sessionToken.Value]
		delete(sessions, sessionToken.Value) // Remove session token
		user := users[username]
		user.SessionToken = ""
		user.CSRFToken = ""
		users[username] = user
	}
	fmt.Fprintln(w, "Logged Out Successfully!")
}

// Exported Authorization Function
var ErrAuth = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	// Retrieve session token from cookies
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Session token not found")
		return ErrAuth
	}
	// Retrieve username from session token
	username, ok := sessions[sessionToken.Value]
	if !ok {
		log.Println("Session token invalid:", sessionToken.Value)
		return ErrAuth
	}
	user, ok := users[username]
	if !ok {
		log.Println("User not found:", username)
		return ErrAuth
	}
	// Validate session token
	if sessionToken.Value != user.SessionToken {
		log.Println("Session token mismatch:", sessionToken.Value, "!=", user.SessionToken)
		return ErrAuth
	}
	// Validate CSRF token
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != user.CSRFToken || csrf == "" {
		log.Println("CSRF token mismatch:", csrf, "!=", user.CSRFToken)
		return ErrAuth
	}
	return nil
}
