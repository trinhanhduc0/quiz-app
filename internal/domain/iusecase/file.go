package usecase

import entity "quiz-app/internal/domain/entities"

type FileUseCase interface {
	CreateFile(file entity.File) (any, error)
	GetFile(file entity.File) (any, error)
	FindByName(file entity.File) (any, error)

	UpdateFile(file entity.File) (any, error)
	DeleteFile(file entity.File) error
}
