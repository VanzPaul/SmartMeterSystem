package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Database environment variables
type DBConfig struct {
	MongoURI string `env:"MONGO_URI"`
	DBName   string `env:"DB_NAME"`
}

var MongoEnv DBConfig

// Logger initialization
var (
	Logger *zap.Logger // Exported logger
	once   sync.Once
)

// InitLogger initializes the logger and ensures it is only done once.
func InitLogger() {
	once.Do(func() {
		var err error

		// Check the DEBUG environment variable
		debug := os.Getenv("DEBUG")
		if debug == "true" {
			// Use Development mode for debugging
			Logger, err = zap.NewDevelopment()
			Logger.Info("Logger initialized in development mode.")
		} else {
			// Use Production mode with custom configuration
			config := zap.NewProductionConfig()
			config.Level.SetLevel(zap.InfoLevel) // Default to Info level
			Logger, err = config.Build()
			Logger.Info("Logger initialized in production mode.")
		}

		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
	})
}

// LoadEnv loads environment variables from the .env file and validates them.
func LoadEnv() error {
	// Load .env file first
	if err := godotenv.Load(); err != nil {
		// Log to stdout since the logger isn't initialized yet
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger after loading environment variables
	InitLogger()

	Logger.Info("Loading environment variables...")

	// Validate required environment variables
	requiredVars := map[string]*string{
		"MONGODB_URI":   &MongoEnv.MongoURI,
		"DATABASE_NAME": &MongoEnv.DBName,
	}
	for key, field := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			Logger.Error("Missing required environment variable", zap.String("key", key))
			return fmt.Errorf("missing required environment variable: %s", key)
		}
		*field = value
	}

	Logger.Info("Environment variables loaded successfully.")
	return nil
}
