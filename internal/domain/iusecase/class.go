package usecase

import (
	entity "quiz-app/internal/domain/entities"
)

type ClassUseCase interface {
	CreateClass(class entity.Class) (any, error)
	UpdateClass(class entity.Class) (any, error)
	DeleteClass(class entity.Class) error
	GetAllClass(id string) ([]any, error)
	GetAllClassByEmail(email string) ([]any, error)
	JoinClass(class entity.Class) error
}
