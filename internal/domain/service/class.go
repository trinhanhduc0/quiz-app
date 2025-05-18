package service

import (
	"context"
	"fmt"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassUseCase struct {
	repoClass repository.ClassRepository
}

func NewClassUseCase(repoClass repository.ClassRepository, repoTest repository.TestRepository) *ClassUseCase {
	return &ClassUseCase{
		repoClass: repoClass,
	}
}

func (uc *ClassUseCase) CreateClass(ctx context.Context, class *entity.Class) (any, error) {
	newID, err := uc.repoClass.CreateClass(ctx, class)
	if err != nil {
		return nil, fmt.Errorf("create class: %w", err)
	}

	class.ID = newID
	return class, nil
}

func (uc *ClassUseCase) GetAllClassByEmail(ctx context.Context, email string) ([]any, error) {
	return uc.repoClass.GetAllClassByEmail(ctx, email)
}

func (uc *ClassUseCase) UpdateClass(ctx context.Context, class *entity.Class) (any, error) {
	return uc.repoClass.UpdateClass(ctx, class)
}

func (uc *ClassUseCase) DeleteClass(ctx context.Context, emailID string, id primitive.ObjectID) error {
	return uc.repoClass.DeleteClass(ctx, emailID, id)
}

func (uc *ClassUseCase) GetAllClass(ctx context.Context, authorEmail string) ([]any, error) {
	return uc.repoClass.GetClassByAuthorEmail(ctx, authorEmail)
}

func (uc *ClassUseCase) JoinClass(ctx context.Context, classID primitive.ObjectID, testIDs []primitive.ObjectID, emailAuthor, studentEmail string) error {

	return uc.repoClass.JoinClass(ctx, classID, studentEmail)
}
func (uc *ClassUseCase) GetAllTestOfClass(ctx context.Context, email string, id primitive.ObjectID) ([]any, error) {
	return uc.repoClass.GetAllTestOfClass(ctx, email, id)
}
func (uc *ClassUseCase) GetQuestionOfTest(ctx context.Context, classId, testId primitive.ObjectID, email string) ([]primitive.ObjectID, primitive.M, error) {
	return uc.repoClass.GetQuestionOfTest(ctx, classId, testId, email)
}
