package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"
	utils "quiz-app/internal/util"
)

// TestMongoRepository implements the repository.TestRepository interface
type TestMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

// NewTestMongoRepository tạo một instance mới của TestMongoRepository
func NewTestMongoRepository() repository.TestRepository {
	collRepo := NewCollRepository("dbapp", "tests")
	return &TestMongoRepository{
		CollRepo: collRepo,
	}
}

// GetTestesByAuthorEmail implements repository.TestRepository.GetTestesByAuthorEmail
func (r *TestMongoRepository) GetTestsByAuthorEmail(ctx context.Context, email string) ([]any, error) {
	filter := bson.M{"author_mail": email}

	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get test by author email: %w", err)
	}

	return results, nil
}

// CreateTest implements repository.TestRepository.CreateTest
func (r *TestMongoRepository) CreateTest(ctx context.Context, class *entity.Test) (primitive.ObjectID, error) {
	classField, err := utils.GenerateUpdateFields(class)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create class: %w", err)
	}
	classField["answer_user"] = bson.A{}
	insertedID, err := r.CollRepo.Create(ctx, classField)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create class: %w", err)
	}

	// Chuyển đổi insertedID về ObjectID
	objID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return objID, nil
}

// UpdateTest implements repository.TestRepository.UpdateTest
func (r *TestMongoRepository) UpdateTest(ctx context.Context, test *entity.Test) (any, error) {
	filter := bson.M{"email_id": test.EmailID, "_id": test.ID}

	testField, err := utils.GenerateUpdateFields(test)

	if err != nil {
		return nil, fmt.Errorf("failed to update test: %w", err)
	}

	results, er := r.CollRepo.Update(ctx, filter, bson.M{"$set": testField})
	if er != nil || results.MatchedCount == 0 {
		return nil, fmt.Errorf("failed to update test: %w", er)
	}
	return test, nil
}

// DeleteTest implements repository.TestRepository.DeleteTest
func (r *TestMongoRepository) DeleteTest(ctx context.Context, id primitive.ObjectID, email string) error {
	filter := bson.M{"email_id": email, "_id": id}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}
	return nil
}
