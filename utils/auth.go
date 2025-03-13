package utils

import (
	"errors"
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
		Logger.Sugar().Debugln("Session token not found")
		return ErrAuth
	}
	Logger.Sugar().Debugf("sessionToken: %s\n", sessionToken)
	// Retrieve username from session token
	email, ok := sessions[sessionToken.Value]
	Logger.Sugar().Debugf("email: %s\n", email)
	if !ok {
		Logger.Sugar().Debugln("Session token invalid:", sessionToken.Value)
		return ErrAuth
	}

	user, ok := users[email]
	Logger.Sugar().Debugf("user: %s\n", user)
	if !ok {
		Logger.Sugar().Debugln("User not found:", email)
		return ErrAuth
	}

	// Validate session token
	if sessionToken.Value != user.SessionToken {
		Logger.Sugar().Debugln("Session token mismatch:", sessionToken.Value, "!=", user.SessionToken)
		return ErrAuth
	}

	// Validate CSRF token: first check header; if missing, use the cookie as a fallback.
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" {
		if csrfCookie, err := r.Cookie("csrf_token"); err == nil {
			csrf = csrfCookie.Value
		}
	}
	if csrf == "" || csrf != user.CSRFToken {
		Logger.Sugar().Debugf("CSRF token mismatch: %s != %s\n", csrf, user.CSRFToken)
		return ErrAuth
	}

	// If everything is valid
	return nil
}
