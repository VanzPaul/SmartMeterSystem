package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSerialNumber = errors.New("invalid serial number")
	secretKey              = []byte("your-secret-key-here") // Replace with a secure secret key
)

// GenerateAPIKey generates a secure API key from a valid serial number
// Serial number must be at least 12 characters long
func GenerateAPIKey(serial string) (string, error) {
	if len(serial) != 12 {
		return "", fmt.Errorf("%w: serial must be 12 characters", ErrInvalidSerialNumber)
	}

	// Create HMAC-SHA256 hash using secret key
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(serial))
	hash := h.Sum(nil)

	// Encode to base64 URL-safe string
	apiKey := base64.URLEncoding.EncodeToString(hash)
	return apiKey, nil
}

// GenerateWifiSSID generates a WiFi SSID from the first 6 characters of the serial number
// with a standardized prefix
func GenerateWifiSSID(serial string) (string, error) {
	if len(serial) < 6 {
		return "", fmt.Errorf("%w: serial must be at least 6 characters", ErrInvalidSerialNumber)
	}

	prefix := "IoT_"
	ssidPart := strings.ToUpper(serial[:6])
	return prefix + ssidPart, nil
}

// GenerateWifiPassword generates a secure WiFi password using the serial number
// Combines hashed value with special characters
func GenerateWifiPassword(serial string) (string, error) {
	if len(serial) < 8 {
		return "", fmt.Errorf("%w: serial must be at least 8 characters", ErrInvalidSerialNumber)
	}

	// Create hash of serial number
	h := sha256.New()
	h.Write([]byte(serial))
	hash := h.Sum(nil)

	// Take first 8 characters of base64 encoded hash and add special characters
	encoded := base64.StdEncoding.EncodeToString(hash)
	password := fmt.Sprintf("%s@%s", encoded[:4], encoded[4:8])
	return password, nil
}
