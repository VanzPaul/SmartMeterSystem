package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vanspaul/SmartMeterSystem/utils"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Database environment variables
type DBConfig struct {
	MongoURI string `env:"MONGO_URI"`
	DBName   string `env:"DB_NAME"`
}

var MongoEnv DBConfig

var (
	once sync.Once
)

// InitLogger initializes the logger from utils and ensures it's done once.
func InitLogger() {
	once.Do(func() {
		debug := os.Getenv("DEBUG") == "true"
		logger, err := utils.NewLogger(debug)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		utils.Logger = logger
	})
}

// LoadEnv loads environment variables and initializes the logger.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	InitLogger()

	utils.Logger.Info("Loading environment variables...")

	requiredVars := map[string]*string{
		"MONGODB_URI":   &MongoEnv.MongoURI,
		"DATABASE_NAME": &MongoEnv.DBName,
	}

	for key, field := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			utils.Logger.Error("Missing required environment variable", zap.String("key", key))
			return fmt.Errorf("missing required environment variable: %s", key)
		}
		*field = value
	}

	utils.Logger.Info("Environment variables loaded successfully.")
	return nil
}
