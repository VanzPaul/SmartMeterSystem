package controllers

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
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

// MongoDBController implements the Database interface for MongoDB.
type MongoDBController struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoDBController creates a new MongoDB controller instance.
func NewMongoDBController(uri, dbName, collName string) (*MongoDBController, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Ping the database to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return &MongoDBController{
		client:     client,
		collection: client.Database(dbName).Collection(collName),
	}, nil
}

// Create inserts a document into the collection.
func (m *MongoDBController) Create(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if document == nil {
		return nil, errors.New("document cannot be nil")
	}
	return m.collection.InsertOne(ctx, document)
}

// Replace updates a document matching the given filter.
func (m *MongoDBController) Replace(ctx context.Context, filter, replacement interface{}) (*mongo.UpdateResult, error) {
	if filter == nil || replacement == nil {
		return nil, errors.New("filter and replacement must be provided")
	}
	return m.collection.ReplaceOne(ctx, filter, replacement)
}

// Update performs a partial update on documents matching the given filter.
func (m *MongoDBController) Update(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error) {
	if filter == nil || update == nil {
		return nil, errors.New("filter and update must be provided")
	}
	return m.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
}

// Delete removes documents matching the given filter.
func (m *MongoDBController) Delete(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if filter == nil {
		return nil, errors.New("filter must be provided")
	}
	return m.collection.DeleteMany(ctx, filter)
}

// Find fetches multiple documents matching the filter.
func (m *MongoDBController) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	if filter == nil {
		filter = bson.D{}
	}
	cursor, err := m.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, results); err != nil {
		return err
	}
	return nil
}

// FindOne fetches a single document matching the filter.
func (m *MongoDBController) FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	if filter == nil {
		filter = bson.D{}
	}
	if result == nil {
		return errors.New("result target cannot be nil")
	}
	singleResult := m.collection.FindOne(ctx, filter, opts...)
	if err := singleResult.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil // No documents found
		}
		return err
	}
	return singleResult.Decode(result)
}

// PaginatedFind supports paginated queries.
func (m *MongoDBController) PaginatedFind(ctx context.Context, filter interface{}, page, perPage int64, results interface{}) error {
	if page < 1 || perPage < 1 {
		return errors.New("page and perPage must be greater than zero")
	}
	opts := options.Find().
		SetSkip((page - 1) * perPage).
		SetLimit(perPage)
	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

// Aggregate executes an aggregation pipeline and returns the results.
func (m *MongoDBController) Aggregate(ctx context.Context, pipeline interface{}, results interface{}) error {
	if pipeline == nil {
		return errors.New("pipeline cannot be nil")
	}
	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, results); err != nil {
		return err
	}
	return nil
}

// WithTransaction supports transactional operations.
func (m *MongoDBController) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, fn)
	return err
}

// Close disconnects the MongoDB client.
func (m *MongoDBController) Close(ctx context.Context) error {
	if m.client == nil {
		return nil // Already closed or not initialized
	}
	return m.client.Disconnect(ctx)
}
