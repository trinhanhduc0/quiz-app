package service

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerUseCase struct {
	repo repository.AnswerRepository
}

func NewAnswerUseCase(repo repository.AnswerRepository) *AnswerUseCase {
	return &AnswerUseCase{repo: repo}
}

func (au *AnswerUseCase) CreateNewAnswer(ctx context.Context, answer *entity.TestAnswer) error {
	newAnswer, err := entity.CreateNewAnswer(answer.TestId, answer.EmailID, answer.Email)
	if err != nil {
		return err
	}

	_, err = au.repo.CreateAnswer(ctx, *newAnswer)
	return err
}

func (au *AnswerUseCase) UpdateAnswer(ctx context.Context, answer entity.TestAnswer) error {
	newAnswer, err := entity.UpdateAnswer(answer)
	if err != nil {
		return err
	}
	fmt.Println("UPDATE")
	fmt.Println(newAnswer)

	_, err = au.repo.UpdateAnswer(ctx, *newAnswer)
	return err
}

func (au *AnswerUseCase) GetAnswer(ctx context.Context, filter primitive.M) (entity.TestAnswer, error) {
	return au.repo.GetAnswer(ctx, filter)
}

func (au *AnswerUseCase) GetAllAnswerByEmail(ctx context.Context, email string) ([]entity.TestAnswer, error) {
	filter := bson.M{"email": email}
	return au.repo.GetAllAnswer(ctx, filter)
}
