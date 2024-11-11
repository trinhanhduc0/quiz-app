package service

import (
	"context"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionUseCase struct {
	QuestionRepo repository.QuestionRepository
}

func NewQuestionUseCase(tr repository.QuestionRepository) *QuestionUseCase {
	return &QuestionUseCase{
		QuestionRepo: tr,
	}
}

// CreateQuestion creates a new question
func (uc *QuestionUseCase) CreateQuestion(ctx context.Context, question *entity.Question) (any, error) {
	newQuestion, err := uc.QuestionRepo.CreateQuestion(ctx, question)
	if err != nil {
		return nil, err
	}
	return newQuestion, nil
}

// GetQuestionByID retrieves a question by its ID
func (uc *QuestionUseCase) GetAllQuestionsByUser(ctx context.Context, userID string, limit, page int) ([]any, error) {
	return uc.QuestionRepo.GetAllQuestionsByUser(ctx, userID, limit, page)
}

// GetQuestionByID retrieves a question by its ID
func (uc *QuestionUseCase) GetAllTestQuestions(ctx context.Context, question_ids []primitive.ObjectID) ([]bson.M, error) {
	return uc.QuestionRepo.GetAllQuestions(ctx, question_ids)
}

// UpdateQuestion updates an existing question
func (uc *QuestionUseCase) UpdateQuestion(ctx context.Context, question *entity.Question) (any, error) {
	return uc.QuestionRepo.UpdateQuestion(ctx, question)
}

// DeleteQuestion deletes a question by ID
func (uc *QuestionUseCase) DeleteQuestion(ctx context.Context, question *entity.Question) error {
	return uc.QuestionRepo.DeleteQuestion(ctx, question)
}
