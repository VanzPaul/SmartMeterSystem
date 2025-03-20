package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

// Cache environment variables
type ServerCacheConfig struct {
	MaxWebUser      string `env:"MAX_WEBUSERS"`
	MaxWebSession   string `env:"MAX_SESSIONS"`
	MaxMeterUser    string `env:"MAX_METERUSERS"`
	MaxMeterSession string `env:"MAX_METERSESSIONS"`
}

var ServerCacheEnv ServerCacheConfig

var (
	MaxWebUsers      int
	MaxWebSessions   int
	MaxMeterUsers    int
	MaxMeterSessions int
)

// Meter evironment variables
type MeterConfig struct {
	SecretKey string `env:"SECRET_KEY"`
}

var MeterEnv MeterConfig

// LoadEnv loads environment variables and initializes the logger.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err) // Fail fast if .env is missing
	}

	// Load required variables FIRST
	// TODO: configure the PORT
	requiredVars := map[string]*string{
		"MONGODB_URI":   &MongoEnv.MongoURI,
		"DATABASE_NAME": &MongoEnv.DBName,
		"DEBUG":         &LogEnv.Debug,
		"SECRET_KEY":    &MeterEnv.SecretKey,
		// "PORT":          ,
		"MAX_WEBUSERS":      &ServerCacheEnv.MaxWebUser,
		"MAX_WEBSESSIONS":   &ServerCacheEnv.MaxWebSession,
		"MAX_METERUSERS":    &ServerCacheEnv.MaxMeterUser,
		"MAX_METERSESSIONS": &ServerCacheEnv.MaxMeterSession,
	}

	for key, field := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("Missing required environment variable: %s", key) // Use standard log
		}
		*field = value
	}

	// Initialize logger AFTER variables are loaded
	initLogger()

	// Initialize cache configuration
	if err := InitializeCacheConfig(); err != nil {
		return fmt.Errorf("failed to initialize cache config: %v", err)
	}

	utils.Logger.Info("Environment variables loaded successfully.")
	return nil
}

// InitLogger uses the now-populated LogEnv.Debug
func initLogger() {
	once.Do(func() {
		debug := LogEnv.Debug == "true"
		logger, err := utils.NewLogger(debug)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		utils.Logger = logger
	})
}

// InitializeCacheConfig converts string values to integers and sets up cache configuration
func InitializeCacheConfig() error {
	var err error
	MaxWebUsers, err = strconv.Atoi(ServerCacheEnv.MaxWebUser)
	if err != nil {
		return fmt.Errorf("invalid MAX_WEBUSERS: %v", err)
	}

	MaxWebSessions, err = strconv.Atoi(ServerCacheEnv.MaxWebSession)
	if err != nil {
		return fmt.Errorf("invalid MAX_SESSIONS: %v", err)
	}

	MaxMeterUsers, err = strconv.Atoi(ServerCacheEnv.MaxMeterUser)
	if err != nil {
		return fmt.Errorf("invalid MAX_METERUSERS: %v", err)
	}

	MaxMeterSessions, err = strconv.Atoi(ServerCacheEnv.MaxMeterSession)
	if err != nil {
		return fmt.Errorf("invalid MAX_METERSESSIONS: %v", err)
	}

	return nil
}
