package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepository interface {
	CreateClass(ctx context.Context, class *entity.Class) (primitive.ObjectID, error)
	GetClassByAuthorEmail(ctx context.Context, email string) ([]any, error)
	UpdateClass(ctx context.Context, class *entity.Class) (any, error)
	DeleteClass(ctx context.Context, emailID string, id primitive.ObjectID) error
	GetAllClassByEmail(ctx context.Context, email string) ([]any, error)
	JoinClass(ctx context.Context, classID primitive.ObjectID, email string) error

	GetAllTestOfClass(ctx context.Context, email string, id primitive.ObjectID) ([]any, error)
	GetQuestionOfTest(ctx context.Context, class, id primitive.ObjectID, email string) ([]primitive.ObjectID, primitive.M, error)
}
