package persistence

import (
	"context"
	"fmt"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollRepository struct {
	Collection *mongo.Collection
}

func NewCollRepository(dbName, collName string) repository.CRUDMongoDB {
	client := GetMongoClient()
	db := client.Database(dbName)
	return &CollRepository{
		Collection: db.Collection(collName),
	}
}

func (r *CollRepository) Create(ctx context.Context, document any) (any, error) {
	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	return result.InsertedID, nil
}

func (r *CollRepository) GetAll(ctx context.Context, filter any) ([]any, error) {
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	var docs []any
	for _, result := range results {
		docs = append(docs, result)
	}
	return docs, nil
}

// GetAllWithOption retrieves documents based on the provided filter and options (e.g., limit, skip)
func (r *CollRepository) GetAllWithOption(ctx context.Context, filter any, opts *options.FindOptions) ([]any, error) {
	// Create an empty slice to store the results
	var results []any

	// Execute the query with the filter and options
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor to decode each document and add to results slice
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, doc)
	}

	// Check if there was an error during iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

func (r *CollRepository) GetWithProjection(ctx context.Context, filter, projection any) ([]any, error) {
	var results []any
	cursor, err := r.Collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return results, nil
}

func (r *CollRepository) GetFilter(ctx context.Context, filter any) (any, error) {
	var result map[string]interface{}
	err := r.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No documents found
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return result, nil
}

func (r *CollRepository) Update(ctx context.Context, filter, update any) (*mongo.UpdateResult, error) {
	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

func (r *CollRepository) UpdateMany(ctx context.Context, filter, update any) (*mongo.UpdateResult, error) {
	result, err := r.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

func (r *CollRepository) Delete(ctx context.Context, filter any) (*mongo.DeleteResult, error) {
	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}
