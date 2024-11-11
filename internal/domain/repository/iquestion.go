package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *entity.Question) (any, error)

	GetAllQuestionsByUser(ctx context.Context, userID string, limit, page int) ([]any, error)

	UpdateQuestion(ctx context.Context, question *entity.Question) (any, error)

	DeleteQuestion(ctx context.Context, question *entity.Question) error

	GetAllQuestions(ctx context.Context, question_ids []primitive.ObjectID) ([]bson.M, error)
}
