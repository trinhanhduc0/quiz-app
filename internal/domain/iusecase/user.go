package usecase

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *entity.User) (primitive.ObjectID, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetUser(ctx context.Context, user *entity.User) *entity.User
}
