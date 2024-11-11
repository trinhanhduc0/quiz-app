package usecase

import (
	"context"
	entity "quiz-app/internal/domain/entities"
)

type QuestionUseCase interface {
	CreateQuestion(ctx context.Context, question *entity.Question) (any, error)
	GetAllQuestionsByUser(ctx context.Context, idUser string) ([]any, error)
	GetAllTestQuestions(ctx context.Context, question_ids []any) ([]any, error)
	UpdateQuestion(ctx context.Context, question *entity.Question) (any, error)
	DeleteQuestion(ctx context.Context, id string) error
}
