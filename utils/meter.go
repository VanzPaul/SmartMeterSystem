package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Package documentation for utility functions related to API key, WiFi SSID, and WiFi password generation.
//
// These functions are designed to generate secure and standardized identifiers based on a serial number.
// They include validation checks to ensure the serial number meets specific requirements and utilize
// cryptographic techniques for security.

var (
	// ErrInvalidSerialNumber is an error returned when the provided serial number does not meet the required format or length.
	ErrInvalidSerialNumber = errors.New("invalid serial number")
)

// GenerateAPIKey generates a secure API key from a valid serial number.
//
// The serial number must be exactly 12 characters long. A cryptographic hash (HMAC-SHA256) is created
// using the secret key, and the resulting hash is encoded into a base64 URL-safe string.
//
// @param serial string - The serial number used to generate the API key. Must be exactly 12 characters.
//
// @return string - The generated API key as a base64 URL-safe string.
// @return error - An error if the serial number is invalid or the hashing process fails.
func GenerateAPIKey(secretKey []byte, serial string) (string, error) {
	// Validate that the serial number is exactly 12 characters long.
	if len(serial) != 12 {
		return "", fmt.Errorf("%w: serial must be 12 characters", ErrInvalidSerialNumber)
	}

	// Create an HMAC-SHA256 hash of the serial number using the secret key.
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(serial))
	hash := h.Sum(nil)

	// Encode the hash into a base64 URL-safe string to create the API key.
	apiKey := base64.URLEncoding.EncodeToString(hash)
	return apiKey, nil
}

// GenerateWifiSSID generates a WiFi SSID from the first 6 characters of the serial number.
//
// The SSID is prefixed with "IoT_" and includes the first 6 characters of the serial number,
// converted to uppercase. The serial number must be at least 6 characters long.
//
// @param serial string - The serial number used to generate the WiFi SSID. Must be at least 6 characters.
//
// @return string - The generated WiFi SSID.
// @return error - An error if the serial number is invalid.
func GenerateWifiSSID(serial string) (string, error) {
	// Validate that the serial number is at least 6 characters long.
	if len(serial) < 12 {
		return "", fmt.Errorf("%w: serial must be at least 12 characters", ErrInvalidSerialNumber)
	}

	// Define the prefix for the SSID and extract the first 6 characters of the serial number.
	prefix := "IoT_"
	ssidPart := strings.ToUpper(serial[:6])

	// Combine the prefix and the extracted part to form the SSID.
	return prefix + ssidPart, nil
}

// GenerateWifiPassword generates a secure WiFi password using the serial number.
//
// The password is derived by hashing the serial number using SHA-256, encoding the hash into base64,
// and formatting it with special characters. The serial number must be at least 8 characters long.
//
// @param serial string - The serial number used to generate the WiFi password. Must be at least 8 characters.
//
// @return string - The generated WiFi password.
// @return error - An error if the serial number is invalid.
func GenerateWifiPassword(serial string) (string, error) {
	// Validate that the serial number is at least 8 characters long.
	if len(serial) < 12 {
		return "", fmt.Errorf("%w: serial must be at least 12 characters", ErrInvalidSerialNumber)
	}

	// Create a SHA-256 hash of the serial number.
	h := sha256.New()
	h.Write([]byte(serial))
	hash := h.Sum(nil)

	// Encode the hash into a base64 string.
	encoded := base64.StdEncoding.EncodeToString(hash)

	// Format the password by taking the first 8 characters of the encoded hash and adding special characters.
	password := fmt.Sprintf("%s@%s", encoded[:4], encoded[4:8])
	return password, nil
}
