package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/models"
)

// Exported Authorization Error
var ErrAuth = errors.New("Unauthorized")

// Authorize function validates session and CSRF tokens
func Authorize(r *http.Request, users map[string]models.LoginData, sessions map[string]string) error {
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
	if csrf == "" || csrf != user.CSRFToken {
		log.Println("CSRF token mismatch:", csrf, "!=", user.CSRFToken)
		return ErrAuth
	}

	// If everything is valid
	return nil
}
