package service

import (
	"context"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
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

func (uc *TestUseCase) GetQuestionOfTest(ctx context.Context, id primitive.ObjectID, email string) ([]primitive.ObjectID, bson.M, error) {
	return uc.TestRepo.GetQuestionOfTest(ctx, id, email)
}

func (cuc *TestUseCase) GetAllTestOfClass(ctx context.Context, email string, id []primitive.ObjectID) ([]any, error) {
	return cuc.TestRepo.GetAllTestOfClass(ctx, email, id)
}

func (cuc *TestUseCase) UpdateAnswerUser(ctx context.Context, testID primitive.ObjectID, email string) error {
	return cuc.TestRepo.UpdateAnswerUser(ctx, testID, email)
}

func (cuc *TestUseCase) AddAllowedUser(ctx context.Context, ids []primitive.ObjectID, user string) error {
	return cuc.TestRepo.AddAllowedUser(ctx, ids, user)
}

func (cuc *TestUseCase) UpdateAllowUser(ctx context.Context, id []primitive.ObjectID, allowedUser []string) error {
	return cuc.TestRepo.UpdateAllowUser(ctx, id, allowedUser)
}

func (cuc *TestUseCase) RemoveAllowUser(ctx context.Context, testID primitive.ObjectID, email string) error {
	return cuc.TestRepo.RemoveAllowUser(ctx, testID, email)
}
