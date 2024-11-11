package usecase

import (
	entity "quiz-app/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
)

type AnswerUseCase interface {
	CreateNewAnswer(answer entity.TestAnswer) error
	UpdateAnswer(answer entity.TestAnswer) error
	GetAnswer(filter bson.M) (entity.TestAnswer, error)
	GetAllAnswer(info bson.M) ([]entity.TestAnswer, error)
}
