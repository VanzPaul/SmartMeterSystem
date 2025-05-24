package database

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func mustStartMongoContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := mongodb.Run(context.Background(), "mongo:latest")
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()
	database = "testdb" // Set test database name

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartMongoContainer()
	if err != nil {
		log.Fatalf("could not start mongodb container: %v", err)
	}

	code := m.Run()

	if err := teardown(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestInsertOne(t *testing.T) {
	srv := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{"name": "Alice"}
	res, err := srv.InsertOne(ctx, "users", doc)
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}
	if res.InsertedID == nil {
		t.Error("InsertedID is nil")
	}

	// Cleanup
	_, _ = srv.DeleteOne(ctx, "users", bson.M{"_id": res.InsertedID})
}

func TestFindOne(t *testing.T) {
	srv := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup
	doc := bson.M{"name": "Bob"}
	insertRes, err := srv.InsertOne(ctx, "users", doc)
	if err != nil {
		t.Fatalf("InsertOne setup failed: %v", err)
	}
	defer srv.DeleteOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})

	// Test FindOne
	result := srv.FindOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})
	var found bson.M
	if err := result.Decode(&found); err != nil {
		t.Fatalf("FindOne decode failed: %v", err)
	}
	if found["name"] != "Bob" {
		t.Errorf("Expected name 'Bob', got %v", found["name"])
	}
}

func TestUpdateOne(t *testing.T) {
	srv := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup
	insertRes, err := srv.InsertOne(ctx, "users", bson.M{"name": "Charlie"})
	if err != nil {
		t.Fatalf("InsertOne setup failed: %v", err)
	}
	defer srv.DeleteOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})

	// Test UpdateOne
	update := bson.M{"$set": bson.M{"name": "Updated"}}
	updateRes, err := srv.UpdateOne(ctx, "users", bson.M{"_id": insertRes.InsertedID}, update)
	if err != nil {
		t.Fatalf("UpdateOne failed: %v", err)
	}
	if updateRes.ModifiedCount != 1 {
		t.Errorf("Expected 1 modification, got %d", updateRes.ModifiedCount)
	}

	// Verify update
	result := srv.FindOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})
	var updated bson.M
	if err := result.Decode(&updated); err != nil {
		t.Fatalf("FindOne after update failed: %v", err)
	}
	if updated["name"] != "Updated" {
		t.Errorf("Expected name 'Updated', got %v", updated["name"])
	}
}

func TestDeleteOne(t *testing.T) {
	srv := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup
	insertRes, err := srv.InsertOne(ctx, "users", bson.M{"name": "Dave"})
	if err != nil {
		t.Fatalf("InsertOne setup failed: %v", err)
	}

	// Test DeleteOne
	deleteRes, err := srv.DeleteOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})
	if err != nil {
		t.Fatalf("DeleteOne failed: %v", err)
	}
	if deleteRes.DeletedCount != 1 {
		t.Errorf("Expected 1 deletion, got %d", deleteRes.DeletedCount)
	}

	// Verify deletion
	result := srv.FindOne(ctx, "users", bson.M{"_id": insertRes.InsertedID})
	if err := result.Err(); err != mongo.ErrNoDocuments {
		t.Errorf("Expected document not found, got %v", err)
	}
}
