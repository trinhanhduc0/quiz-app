package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestRepository interface {
	CreateTest(ctx context.Context, test *entity.Test) (primitive.ObjectID, error)
	GetTestsByAuthorEmail(ctx context.Context, email string) ([]any, error)
	UpdateTest(ctx context.Context, test *entity.Test) (any, error)
	DeleteTest(ctx context.Context, id primitive.ObjectID, email string) error
	
}
