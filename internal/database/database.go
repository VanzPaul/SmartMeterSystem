package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	InsertOne(context.Context, string, interface{}) (*mongo.InsertOneResult, error)
	FindOne(context.Context, string, interface{}) *mongo.SingleResult
	FindMany(context.Context, string, interface{}) (*mongo.Cursor, error)
	UpdateOne(context.Context, string, interface{}, interface{}) (*mongo.UpdateResult, error)
	DeleteOne(context.Context, string, interface{}) (*mongo.DeleteResult, error)
	Aggregation(context.Context, string, interface{}) (*mongo.Cursor, error)
}

type service struct {
	db *mongo.Database
}

var (
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	database = os.Getenv("BLUEPRINT_DB_DATABASE")
	username = os.Getenv("BLUEPRINT_DB_USERNAME")
	password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
)

func New() Service {
	if host == "" {
		log.Fatalln("host is empty")
	} else if port == "" {
		log.Fatalln("port is empty")
	} else if database == "" {
		log.Fatalln("databse is empty")
	} else if username == "" {
		log.Fatalln("username is empty")
	} else if password == "" {
		log.Fatal("password is empty")
	}

	// Build connection URI conditionally
	// mongodb://melkey:password1234@localhost:27017/
	var uri string
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/",
			url.QueryEscape(username),
			url.QueryEscape(password),
			host,
			port,
		)
		log.Println(uri)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/%s/",
			host,
			port,
			database,
		)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the database name is provided
	if database == "" {
		log.Fatal("Database name must be provided")
	}

	return &service{
		db: client.Database(database),
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database using the client from the database instance
	err := s.db.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// InsertOne inserts a single document into the specified collection
func (s *service) InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := s.db.Collection(collection)
	return coll.InsertOne(ctx, document)
}

// FindOne finds a single document in the specified collection based on the filter
func (s *service) FindOne(ctx context.Context, collection string, filter interface{}) *mongo.SingleResult {
	coll := s.db.Collection(collection)
	return coll.FindOne(ctx, filter)
}

func (s *service) FindMany(ctx context.Context, collection string, filter interface{}) (*mongo.Cursor, error) {
	coll := s.db.Collection(collection)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// UpdateOne updates a single document in the specified collection based on the filter
func (s *service) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	coll := s.db.Collection(collection)
	return coll.UpdateOne(ctx, filter, update)
}

// DeleteOne deletes a single document from the specified collection based on the filter
func (s *service) DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := s.db.Collection(collection)
	return coll.DeleteOne(ctx, filter)
}

// Aggregation performs an aggregation operation on the specified collection.
func (s *service) Aggregation(ctx context.Context, collection string, pipeline interface{}) (*mongo.Cursor, error) {
	coll := s.db.Collection(collection)
	return coll.Aggregate(ctx, pipeline)
}
