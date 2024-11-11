package usecase

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestUseCase interface {
	CreateTest(ctx context.Context, test *entity.Test) (primitive.ObjectID, error)
	GetTestByID(ctx context.Context, id primitive.ObjectID) (any, error)
	GetTestsByAuthorEmail(ctx context.Context, email string) ([]any, error)
	UpdateTest(ctx context.Context, test *entity.Test) error
	DeleteTest(ctx context.Context, id primitive.ObjectID) error
	GetAllTestOfClass(email string, id []primitive.ObjectID) ([]any, error)
	GetQuestionOfTest(ctx context.Context, question_ids any) ([]any, bson.M, error)

	UpdateAllowUser(ctx context.Context, id []primitive.ObjectID, allowedUser []string) error
	UpdateAnswerUser(ctx context.Context, testID primitive.ObjectID, email string) error
}
