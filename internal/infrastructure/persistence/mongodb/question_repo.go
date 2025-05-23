package persistence

import (
	"context"
	"fmt"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QuestionMongoRepository implements the repository.QuestionRepository interface
type QuestionMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

// NewQuestionMongoRepository creates a new instance of QuestionMongoRepository
func NewQuestionMongoRepository() repository.QuestionRepository {
	collRepo := NewCollRepository("dbapp", "questions")
	return &QuestionMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *QuestionMongoRepository) GetAllQuestions(ctx context.Context, questionIDs []primitive.ObjectID) ([]primitive.M, error) {
	if len(questionIDs) == 0 {
		return nil, fmt.Errorf("no question IDs provided")
	}

	// Lọc theo danh sách ID
	questionFilter := bson.M{"_id": bson.M{"$in": questionIDs}}

	// Projection để loại bỏ các trường không cần
	projection := bson.M{
		"createdAt": 0,
		"updatedAt": 0,
		"metadata":  0,
	}

	// Fetch với projection
	questionList, err := r.CollRepo.GetWithProjection(ctx, questionFilter, projection)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions with filter %v: %w", questionFilter, err)
	}

	fmt.Println(questionList)
	var results []primitive.M
	for _, item := range questionList {
		fmt.Printf("%T", item)
		if question, ok := item.(primitive.M); ok {
			primitiveM := primitive.M(question)
			results = append(results, primitiveM)
		} else {
			fmt.Printf("unexpected type in questionList: %T", item)
			return nil, fmt.Errorf("unexpected type in questionList: %T", item)
		}
	}

	return results, nil
}

// CreateQuestion implements repository.QuestionRepository.CreateQuestion
func (r *QuestionMongoRepository) CreateQuestion(ctx context.Context, question *entity.Question) (any, error) {
	insertedID, err := r.CollRepo.Create(ctx, question)
	if err != nil {
		return question, fmt.Errorf("failed to create question: %w", err)
	}

	question.ID = insertedID.(primitive.ObjectID)
	return question, nil
}

func (r *QuestionMongoRepository) GetAllQuestionsByUser(ctx context.Context, email_id string, limit, page int) ([]any, error) {
	// Create a filter to get questions by email_id
	filter := bson.M{"metadata.author": email_id}
	fmt.Println("AUTHOR: ", email_id)
	// Set options for pagination
	findOpts := options.Find()
	findOpts.SetLimit(int64(limit))
	findOpts.SetSkip(int64(page * limit)) // Calculate offset based on page and limit

	// Fetch questions using the filter and options
	allQuestions, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions by email_id: %w", err)
	}
	fmt.Println("QUESTION: ", allQuestions)
	return allQuestions, nil
}

// UpdateQuestion implements repository.QuestionRepository.UpdateQuestion
func (r *QuestionMongoRepository) UpdateQuestion(ctx context.Context, question *entity.Question) (any, error) {
	filter := bson.M{"metadata.author": question.Metadata.Author, "_id": question.ID}

	result, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": question})
	fmt.Println(result)
	if err != nil {
		return &entity.Question{}, fmt.Errorf("failed to update question: %w", err)
	}
	// questionField["_id"] = question.ID
	return question, nil
}

// DeleteQuestion implements repository.QuestionRepository.DeleteQuestion
func (r *QuestionMongoRepository) DeleteQuestion(ctx context.Context, question *entity.Question) error {
	filter := bson.M{"metadata.author": question.Metadata.Author, "_id": question.ID}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	return nil
}
