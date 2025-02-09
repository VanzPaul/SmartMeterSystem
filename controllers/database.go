package controllers

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database interface defines the contract for database operations.
type Database interface {
	Create(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	Replace(ctx context.Context, filter, replacement interface{}) (*mongo.UpdateResult, error)
	Update(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error
	FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error
	Aggregate(ctx context.Context, pipeline interface{}, results interface{}) error // New method for aggregation
	Close(ctx context.Context) error
}

type MongoDB struct {
	client *mongo.Client
	dbName string
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &MongoDB{
		client: client,
		dbName: dbName,
	}, nil
}

func (db *MongoDB) Create(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	coll := db.client.Database(db.dbName).Collection("default")
	return coll.InsertOne(ctx, document)
}

func (db *MongoDB) Replace(ctx context.Context, filter, replacement interface{}) (*mongo.UpdateResult, error) {
	coll := db.client.Database(db.dbName).Collection("default")
	return coll.ReplaceOne(ctx, filter, replacement)
}

func (db *MongoDB) Update(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := db.client.Database(db.dbName).Collection("default")
	return coll.UpdateOne(ctx, filter, update)
}

func (db *MongoDB) Delete(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	coll := db.client.Database(db.dbName).Collection("default")
	return coll.DeleteOne(ctx, filter)
}

func (db *MongoDB) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	coll := db.client.Database(db.dbName).Collection("default")
	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

func (db *MongoDB) FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	coll := db.client.Database(db.dbName).Collection("default")
	return coll.FindOne(ctx, filter, opts...).Decode(result)
}

func (db *MongoDB) Aggregate(ctx context.Context, pipeline interface{}, results interface{}) error {
	coll := db.client.Database(db.dbName).Collection("default")
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
