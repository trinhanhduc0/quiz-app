package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
)

type AnswerRepository interface {
	CreateAnswer(ctx context.Context, answer entity.TestAnswer) (*entity.TestAnswer, error)
	UpdateAnswer(ctx context.Context, answer entity.TestAnswer) (*entity.TestAnswer, error)
	GetAnswer(ctx context.Context, filter bson.M) (entity.TestAnswer, error)
	GetAllAnswer(ctx context.Context, infoAnswer bson.M) ([]entity.TestAnswer, error)
}
