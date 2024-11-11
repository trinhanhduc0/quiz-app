package persistence

import (
	"context"
	"fmt"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"
	utils "quiz-app/internal/util"

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
	// Kiểm tra xem questionIDs có rỗng không
	if len(questionIDs) == 0 {
		return nil, fmt.Errorf("no question IDs provided")
	}

	// Sử dụng $in để lọc câu hỏi theo danh sách ObjectIDs
	questionFilter := bson.M{"_id": bson.M{"$in": questionIDs}}

	// Fetch all matching questions
	questionList, err := r.CollRepo.GetAll(ctx, questionFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions from test with filter %v: %w", questionFilter, err)
	}

	// Chuyển đổi questionList sang []primitive.M
	var results []primitive.M
	for _, item := range questionList {
		// Kiểm tra kiểu của item
		if question, ok := item.(map[string]interface{}); ok {
			// Chuyển đổi map[string]interface{} sang primitive.M
			primitiveM := primitive.M(question) // ép kiểu
			results = append(results, primitiveM)
		} else {
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

func (r *QuestionMongoRepository) GetAllQuestionsByUser(ctx context.Context, userID string, limit, page int) ([]any, error) {
	// Create a filter to get questions by userID
	filter := bson.M{"metadata.author": userID}

	// Set options for pagination
	findOpts := options.Find()
	findOpts.SetLimit(int64(limit))
	findOpts.SetSkip(int64(page * limit)) // Calculate offset based on page and limit

	// Fetch questions using the filter and options
	allQuestions, err := r.CollRepo.GetAllWithOption(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions by userID: %w", err)
	}

	return allQuestions, nil
}

// UpdateQuestion implements repository.QuestionRepository.UpdateQuestion
func (r *QuestionMongoRepository) UpdateQuestion(ctx context.Context, question *entity.Question) (any, error) {
	filter := bson.M{"metadata.author": question.Metadata.Author, "_id": question.ID}

	questionField, err := utils.GenerateUpdateFields(question)

	if err != nil {
		return &entity.Question{}, fmt.Errorf("failed to update question: %w", err)
	}

	_, err = r.CollRepo.Update(ctx, filter, bson.M{"$set": questionField})
	if err != nil {
		return &entity.Question{}, fmt.Errorf("failed to update question: %w", err)
	}
	questionField["_id"] = question.ID
	return questionField, nil
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
