package service

import (
	"context"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestUseCase struct {
	TestRepo repository.TestRepository
}

func NewTestUseCase(tr repository.TestRepository) *TestUseCase {
	return &TestUseCase{
		TestRepo: tr,
	}
}

func (uc *TestUseCase) CreateTest(ctx context.Context, test *entity.Test) (primitive.ObjectID, error) {
	return uc.TestRepo.CreateTest(ctx, test)
}

func (uc *TestUseCase) GetTestsByAuthorEmail(ctx context.Context, email string) ([]any, error) {
	return uc.TestRepo.GetTestsByAuthorEmail(ctx, email)
}

func (uc *TestUseCase) UpdateTest(ctx context.Context, test *entity.Test) (any, error) {
	return uc.TestRepo.UpdateTest(ctx, test)
}

func (uc *TestUseCase) DeleteTest(ctx context.Context, id primitive.ObjectID, email string) error {
	return uc.TestRepo.DeleteTest(ctx, id, email)
}



