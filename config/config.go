package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnvLoader defines an interface for loading environment variables.
type EnvLoader interface {
	GetEnv(key string) string
}

// DefaultEnvLoader implements the EnvLoader interface using the actual environment.
type DefaultEnvLoader struct{}

func (d DefaultEnvLoader) GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnv retrieves an environment variable using the DefaultEnvLoader.
func GetEnv(key string) string {
	loader := DefaultEnvLoader{}
	return loader.GetEnv(key)
}

// LoadEnv loads environment variables from the .env file.
func LoadEnv() {
	fmt.Println("Loading .env file...")
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	fmt.Println(".env file loaded successfully.")
}

// ConnectDB connects to MongoDB using the provided EnvLoader.
func ConnectDB(loader EnvLoader) *mongo.Client {
	fmt.Println("Connecting to MongoDB...")

	// Retrieve the MongoDB URI from the environment variables
	uri := loader.GetEnv("MONGODB_URI")
	fmt.Println("MongoDB URI retrieved from environment:", uri)

	if uri == "" {
		log.Fatal("MONGODB_URI is not set in the .env file.")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	fmt.Println("Attempting to connect to MongoDB with the following options:")
	fmt.Println("URI:", uri)
	fmt.Println("Server API Version:", serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	fmt.Println("Successfully connected to MongoDB.")

	// Send a ping to confirm a successful connection
	fmt.Println("Pinging MongoDB deployment...")
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}
	fmt.Println("Ping successful. MongoDB connection is active.")

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}
