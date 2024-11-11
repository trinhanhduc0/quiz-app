package persistence

import (
	"context"
	"errors"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnswerMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

func NewAnswerMongoRepository() repository.AnswerRepository {
	collRepo := NewCollRepository("dbapp", "answers")
	return &AnswerMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *AnswerMongoRepository) CreateAnswer(ctx context.Context, answer entity.TestAnswer) (*entity.TestAnswer, error) {
	insertedID, err := r.CollRepo.Create(ctx, answer)
	if err != nil {
		return nil, err
	}

	objID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}

	answer.ID = objID
	return &answer, nil
}

func (r *AnswerMongoRepository) UpdateAnswer(ctx context.Context, answer entity.TestAnswer) (*entity.TestAnswer, error) {
	filter := primitive.M{"test_id": answer.TestId, "email_id": answer.EmailID}

	updateResult, err := r.CollRepo.Update(ctx, filter, primitive.M{"$set": answer})
	if err != nil || updateResult.MatchedCount == 0 {
		return nil, err
	}

	return &answer, nil
}

func (r *AnswerMongoRepository) GetAnswer(ctx context.Context, filter bson.M) (entity.TestAnswer, error) {
	result, err := r.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.TestAnswer{}, errors.New("no answer found")
		}
		return entity.TestAnswer{}, fmt.Errorf("error retrieving answer: %v", err)
	}

	var answer entity.TestAnswer
	bsonBytes, err := bson.Marshal(result) // Marshal the result
	if err != nil {
		return entity.TestAnswer{}, fmt.Errorf("error marshaling result: %v", err)
	}

	err = bson.Unmarshal(bsonBytes, &answer) // Unmarshal into answer
	if err != nil {
		return entity.TestAnswer{}, fmt.Errorf("error unmarshaling result: %v", err)
	}

	return answer, nil
}
func (r *AnswerMongoRepository) GetAllAnswer(ctx context.Context, filter bson.M) ([]entity.TestAnswer, error) {
	projection := bson.M{}
	results, err := r.CollRepo.GetWithProjection(ctx, filter, projection)
	if err != nil {
		return nil, err
	}

	var answers []entity.TestAnswer
	for _, result := range results {
		var answer entity.TestAnswer
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &answer)
		answers = append(answers, answer)
	}
	return answers, nil
}
