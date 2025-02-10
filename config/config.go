package config

import (
	"log"
	"os"
	"sync"

	"github.com/vanspaul/SmartMeterSystem/utils"

	"github.com/joho/godotenv"
)

// Database environment variables
type DBConfig struct {
	MongoURI string `env:"MONGO_URI"`
	DBName   string `env:"DB_NAME"`
}

var MongoEnv DBConfig

type LogConfig struct {
	Debug string `env:"DEBUG"`
}

var LogEnv LogConfig

var (
	once sync.Once
)

// LoadEnv loads environment variables and initializes the logger.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err) // Fail fast if .env is missing
	}

	// Load required variables FIRST
	requiredVars := map[string]*string{
		"MONGODB_URI":   &MongoEnv.MongoURI,
		"DATABASE_NAME": &MongoEnv.DBName,
		"DEBUG":         &LogEnv.Debug,
	}

	for key, field := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("Missing required environment variable: %s", key) // Use standard log
		}
		*field = value
	}

	// Initialize logger AFTER variables are loaded
	InitLogger()

	utils.Logger.Info("Environment variables loaded successfully.")
	return nil
}

// InitLogger uses the now-populated LogEnv.Debug
func InitLogger() {
	once.Do(func() {
		debug := LogEnv.Debug == "true"
		logger, err := utils.NewLogger(debug)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		utils.Logger = logger
	})
}
