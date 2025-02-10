package controllers

import (
	"context"
	"fmt"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Database interface defines the contract for database operations.
type Database interface {
	Create(ctx context.Context, collname models.Collection, document interface{}) (*mongo.InsertOneResult, error)
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

type MongoEnv string

var (
	MongoURI MongoEnv
	DBName   MongoEnv
)

func NewMongoDB(ctx context.Context, env *config.DBConfig) (*MongoDB, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(env.MongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	config.Logger.Debug("Connected to MongoDB Succesfully", zap.String("URI", env.MongoURI), zap.String("DBName", env.DBName))

	// Ping the database to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		config.Logger.Error("failed to ping MongoDB", zap.String("URI", env.MongoURI), zap.String("DBName", env.DBName))
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	config.Logger.Debug("Pinged MongoDB Succesfully")

	return &MongoDB{
		client: client,
		dbName: env.DBName,
	}, nil
}

func (db *MongoDB) Create(ctx context.Context, collname models.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	return coll.InsertOne(ctx, document)
}

func (db *MongoDB) Replace(ctx context.Context, collname models.Collection, filter, replacement interface{}) (*mongo.UpdateResult, error) {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	return coll.ReplaceOne(ctx, filter, replacement)
}

func (db *MongoDB) Update(ctx context.Context, collname models.Collection, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	return coll.UpdateOne(ctx, filter, update)
}

func (db *MongoDB) Delete(ctx context.Context, collname models.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	return coll.DeleteOne(ctx, filter)
}

func (db *MongoDB) Find(ctx context.Context, collname models.Collection, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

func (db *MongoDB) FindOne(ctx context.Context, collname models.Collection, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	coll := db.client.Database(db.dbName).Collection(string(collname))
	return coll.FindOne(ctx, filter, opts...).Decode(result)
}

func (db *MongoDB) Aggregate(ctx context.Context, collname models.Collection, pipeline interface{}, results interface{}) error {
	coll := db.client.Database(db.dbName).Collection(string(collname))
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
