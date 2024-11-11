package service

import (
	"context"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(tr repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: tr,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *entity.User) (primitive.ObjectID, error) {
	return uc.UserRepo.CreateUser(ctx, user)
}

func (uc *UserUseCase) GetUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return uc.UserRepo.GetUser(ctx, user)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return uc.UserRepo.UpdateUser(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return uc.UserRepo.DeleteUser(ctx, id)
}
