package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CRUDMongoDB interface {
	Create(ctx context.Context, document any) (any, error)
	GetAll(ctx context.Context, filter any) ([]any, error)
	GetAllWithOption(ctx context.Context, filter any, opts *options.FindOptions) ([]any, error)
	GetFilter(ctx context.Context, filter any) (any, error)
	GetWithProjection(ctx context.Context, filter, projection any) ([]any, error)
	Update(ctx context.Context, filter, update any) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter, update any) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, filter any) (*mongo.DeleteResult, error)
}
