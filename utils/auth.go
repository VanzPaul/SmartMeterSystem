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
	username, ok := sessions[sessionToken.Value]
	Logger.Sugar().Debugf("username: %s\n", username)
	if !ok {
		Logger.Sugar().Debugln("Session token invalid:", sessionToken.Value)
		return ErrAuth
	}

	user, ok := users[username]
	Logger.Sugar().Debugf("user: %s\n", user)
	if !ok {
		Logger.Sugar().Debugln("User not found:", username)
		return ErrAuth
	}

	// Validate session token
	if sessionToken.Value != user.SessionToken {
		Logger.Sugar().Debugln("Session token mismatch:", sessionToken.Value, "!=", user.SessionToken)
		return ErrAuth
	}

	// Validate CSRF token
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" || csrf != user.CSRFToken {
		Logger.Sugar().Debugln("CSRF token mismatch:", csrf, "!=", user.CSRFToken)
		return ErrAuth
	}

	// If everything is valid
	return nil
}
