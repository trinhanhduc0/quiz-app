package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestRepository interface {
	CreateTest(ctx context.Context, test *entity.Test) (primitive.ObjectID, error)
	GetTestsByAuthorEmail(ctx context.Context, email string) ([]any, error)
	UpdateTest(ctx context.Context, test *entity.Test) (any, error)
	DeleteTest(ctx context.Context, id primitive.ObjectID, email string) error
	GetAllTestOfClass(ctx context.Context, email string, id []primitive.ObjectID) ([]any, error)
	GetQuestionOfTest(ctx context.Context, id primitive.ObjectID, email string) ([]primitive.ObjectID, bson.M, error)
	
	UpdateAllowUser(ctx context.Context, id []primitive.ObjectID, allowedUser []string) error
	AddAllowedUser(ctx context.Context, ids []primitive.ObjectID, user string) error
	UpdateAnswerUser(ctx context.Context, testID primitive.ObjectID, email string) error
	RemoveAllowUser(ctx context.Context, testID primitive.ObjectID, email string) error
}
