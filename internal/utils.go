package internal

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create logger
	isDevelopment, parseErr := strconv.ParseBool(os.Getenv("DEBUG"))
	if parseErr != nil {
		panic(parseErr)
	}

	if !isDevelopment {
		// Use NewProduction for production environments, which includes sensible defaults like JSON format and log rotation.
		logger, err = zap.NewProduction()
	} else {
		// Use NewDevelopment for development environments, which provides more human-readable output.
		// Create a console encoder for terminal output
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

		// Create a JSON encoder for file output
		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

		// Define outputs
		consoleOutput := zapcore.Lock(os.Stdout) // Terminal output
		fileOutput, openFileErr := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if openFileErr != nil {
			panic(openFileErr)
		}

		// Combine outputs with different encoders
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleOutput, zapcore.DebugLevel),
			zapcore.NewCore(fileEncoder, zapcore.AddSync(fileOutput), zapcore.InfoLevel),
		)

		// Build the logger
		logger = zap.New(core)
	}

	if err != nil {
		return nil, err // Return the error if logger creation fails
	}

	return logger, nil
}
