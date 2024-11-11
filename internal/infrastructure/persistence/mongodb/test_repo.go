package persistence

import (
	"context"
	"fmt"
	"time"

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
	filter := bson.M{"email_id": email}

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

func (r *TestMongoRepository) GetQuestionOfTest(ctx context.Context, id primitive.ObjectID, email string) ([]primitive.ObjectID, bson.M, error) {
	// Use $in to check if the user's email is in the allowed_users array
	filter := bson.M{
		"_id": id,
		"$or": []bson.M{
			{"allowed_users": bson.M{"$in": []string{email}}},
			{"answer_user": bson.M{"$in": []string{email}}},
		},
	}
	projection := bson.M{"question_ids": 1, "start_time": 1, "end_time": 1, "allowed_users": 1, "updated_time": 1, "random": 1, "is_test": 1, "duration_minutes": 1, "_id": 0}

	//Infor question
	var questionInfo bson.M
	var listQuestions []any
	var err error
	// Get all matching results
	listQuestions, err = r.CollRepo.GetWithProjection(ctx, filter, projection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get question IDs for the test: %w", err)
	}
	// Ensure that listQuestions is of type []bson.M
	var questionIDs []primitive.ObjectID

	for i := 0; i < len(listQuestions); i++ {
		// Type assert the document as bson.M and extract question_ids
		doc, ok := listQuestions[i].(bson.M)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected type for document; expected bson.M")
		}

		questions, ok := doc["question_ids"].(bson.A)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected type for question_ids field; expected bson.A (array)")
		}

		// Convert the []interface{} to []primitive.ObjectID
		for _, q := range questions {
			questionID, ok := q.(primitive.ObjectID)
			if !ok {
				return nil, nil, fmt.Errorf("invalid question ID format; expected ObjectID")
			}
			questionIDs = append(questionIDs, questionID)
		}
		questionInfo = doc
		delete(listQuestions[i].(bson.M), "question_ids")
	}

	return questionIDs, questionInfo, nil
}

// GetAllTestOfClass retrieves all tests for a class by student email and class IDs.
func (r *TestMongoRepository) GetAllTestOfClass(ctx context.Context, email string, ids []primitive.ObjectID) ([]any, error) {
	// Define the filter to match classes where the student is accepted and class ID matches
	filter := bson.M{"_id": bson.M{"$in": ids}}

	// Specify the fields you want to include in the result
	projection :=
		bson.M{
			"class_name":       1,
			"descript":         1,
			"test_name":        1,
			"tags":             1,
			"duration_minutes": 1,
			"_id":              1,
			"start_time":       1,
			"end_time":         1,
			"is_test":          1,
		}

	// Get the matching documents with the specified projection
	classes, err := r.CollRepo.GetWithProjection(ctx, filter, projection)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes by email: %w", err)
	}

	return classes, nil
}
func (r *TestMongoRepository) UpdateAllowUser(ctx context.Context, ids []primitive.ObjectID, allowedUsers []string) error {
	// Define the filter to match documents by the provided ids
	filter := bson.M{"_id": bson.M{"$in": ids}}

	// Define the update to set the allowed_users field
	update := bson.M{
		"$set": bson.M{
			"allowed_users": allowedUsers,
			"updated_at":    time.Now(), // Optionally update the updated_at field
		},
	}

	// Perform the update operation
	result, err := r.CollRepo.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update allowed users: %v", err)
	}

	fmt.Printf("Matched %d documents and modified %d documents.\n", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *TestMongoRepository) RemoveAllowUser(ctx context.Context, testID primitive.ObjectID, email string) error {
	filter := bson.M{"_id": testID}

	// Toán tử $pull sẽ xóa email khỏi mảng nếu tồn tại.
	update := bson.M{
		"$pull": bson.M{
			"allowed_users": email,
		},
	}

	// Thực hiện cập nhật
	result, err := r.CollRepo.Update(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove answer_user: %v", err)
	}

	// Nếu không tìm thấy tài liệu với ID tương ứng
	if result.MatchedCount == 0 {
		return fmt.Errorf("no test found with the given ID")
	}

	return nil
}

func (r *TestMongoRepository) UpdateAnswerUser(ctx context.Context, testID primitive.ObjectID, email string) error {
	filter := bson.M{"_id": testID}

	// Toán tử $addToSet sẽ thêm email vào mảng nếu chưa tồn tại.
	update := bson.M{
		"$addToSet": bson.M{
			"answer_user": email,
		},
	}

	// Thực hiện cập nhật
	result, err := r.CollRepo.Update(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update answer_user: %v", err)
	}

	// Nếu không tìm thấy tài liệu với ID tương ứng
	if result.MatchedCount == 0 {
		return fmt.Errorf("no test found with the given ID")
	}

	return nil
}
