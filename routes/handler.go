package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vanspaul/SmartMeterSystem/utils"
)

type login struct {
	HashedPassword, SessionToken, CSRFToken string
}

// Database. change this with mongodb
var users = map[string]login{}

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		er := http.StatusNotAcceptable
		http.Error(w, "Invalid username/password", er)
		return
	}

	if _, ok := users[username]; ok {
		er := http.StatusConflict
		http.Error(w, "User already exist", er)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	users[username] = login{
		HashedPassword: hashedPassword,
	}

	fmt.Fprintln(w, "User Registered Succesfulluy")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !utils.CheckPasswordHash(password, user.HashedPassword) {
		er := http.StatusUnauthorized
		http.Error(w, "Invalid Username or Password", er)
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

	// Store session & csrf cookie in the database
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	fmt.Fprintln(w, "Login Successful!")
}

func Protected(w http.ResponseWriter, r *http.Request) {
	if err := authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized", er)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s", username)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if err := authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthoried", er)
		return
	}

	// Clear cookie
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
	username := r.FormValue("username")
	user := users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	users[username] = user

	fmt.Fprintln(w, "Logged Out Succesfully!")
}

// Sessions
var ErrAuth = errors.New("Unauthorized")

func authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, ok := users[username]
	if !ok {
		return ErrAuth
	}

	// Get the Sessiob Token from the cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value != user.SessionToken {
		return ErrAuth
	}

	// Get the CSRF from the headers
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != user.CSRFToken || csrf == "" {
		return ErrAuth
	}

	return nil
}
