package internal

import (
	"net"
	"os"
	"strconv"

	"go.uber.org/zap"
)

func GetResolvedIP() string {
	// Print the address where the server is running
	// Resolve the hostname and IP address
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	var resolvedIP string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			resolvedIP = ipv4.String()
			break
		}
	}
	return resolvedIP
}

// InitLogger initializes a zap logger based on the environment (e.g., production or development)
func NewLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	// Create logger
	isDevelopment, parseErr := strconv.ParseBool(os.Getenv("DEBUG"))
	if parseErr != nil {
		panic(parseErr)
	}

	if !isDevelopment {
		// Use NewProduction for production environments, which includes sensible defaults like JSON format and log rotation.
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err // Return the error if logger creation fails
	}

	return logger, nil
}
