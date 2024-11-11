package repository

import (
	"context"
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileRepository interface {
	CreateFile(ctx context.Context, file *entity.File) (primitive.ObjectID, error)
	GetFile(ctx context.Context, file *entity.File) (any, error)
	FindByName(ctx context.Context, file *entity.File) (any, error)

	UpdateFile(ctx context.Context, file *entity.File) (any, error)
	DeleteFile(ctx context.Context, file entity.File) error

	GetAllImageFile(ctx context.Context, email string) ([]any, error)
	GetAllFile(ctx context.Context, email string) (any, error)
}
