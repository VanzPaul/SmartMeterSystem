package config

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockEnvLoader is a mock implementation of the EnvLoader interface.
type MockEnvLoader struct {
	mockEnv map[string]string
}

func (m MockEnvLoader) GetEnv(key string) string {
	return m.mockEnv[key]
}

// TestGetEnv tests the GetEnv function.
func TestGetEnv(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	// Test case 1: Valid key
	value := GetEnv("TEST_KEY")
	assert.Equal(t, "test_value", value, "Expected value to be 'test_value'")

	// Test case 2: Invalid key
	value = GetEnv("INVALID_KEY")
	assert.Empty(t, value, "Expected value to be empty for invalid key")
}

// TestConnectDB tests the ConnectDB function.
func TestConnectDB(t *testing.T) {
	// Define test credentials and URI
	username := "database_test"
	password := "password"
	ip := "192.168.1.12"
	port := "27017"
	authMechanism := "SCRAM-SHA-256"

	// Construct the MongoDB URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authMechanism=%s", username, password, ip, port, authMechanism)

	// Mock environment variables
	mockLoader := MockEnvLoader{
		mockEnv: map[string]string{
			"MONGODB_URI":      uri,
			"MONGODB_USERNAME": username,
			"MONGODB_PASSWORD": password,
		},
	}

	// Mock MongoDB client
	mockClient := &mongo.Client{}
	mockConnect := func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
		return mockClient, nil
	}

	// Override the mongo.Connect function
	oldConnect := mongoConnect
	defer func() { mongoConnect = oldConnect }()
	mongoConnect = mockConnect

	// Call ConnectDB with the mock loader
	client := ConnectDB(mockLoader)

	// Assertions
	assert.NotNil(t, client, "Expected a non-nil MongoDB client")
	assert.Equal(t, mockClient, client, "Expected the mock MongoDB client to be returned")
}

// mongoConnect is a variable to hold the mongo.Connect function for mocking.
var mongoConnect = mongo.Connect
