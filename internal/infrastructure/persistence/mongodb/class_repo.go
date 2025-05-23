package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"
)

// ClassMongoRepository implements the repository.ClassRepository interface
type ClassMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

// NewClassMongoRepository tạo một instance mới của ClassMongoRepository
func NewClassMongoRepository() repository.ClassRepository {
	collRepo := NewCollRepository("dbapp", "classes")
	return &ClassMongoRepository{
		CollRepo: collRepo,
	}
}

// CreateClass implements repository.ClassRepository.CreateClass
func (r *ClassMongoRepository) CreateClass(ctx context.Context, class *entity.Class) (primitive.ObjectID, error) {
	insertedID, err := r.CollRepo.Create(ctx, class)
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

// GetClassByAuthorEmail fetches all classes for a given author's email ID.
func (r *ClassMongoRepository) GetClassByAuthorEmail(ctx context.Context, email string) ([]any, error) {
	// Define a filter based on email_id
	filter := bson.M{"author_mail": email}

	// Query the database to retrieve all matching documents
	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes by author email: %w", err)
	}

	// Return the list of classes
	return results, nil
}

func (r *ClassMongoRepository) GetAllClassByEmail(ctx context.Context, email string) ([]any, error) {
	filter := bson.M{"students_accept": email}
	projection := bson.M{"test_id": 1, "class_name": 1, "author_mail": 1, "tags": 1, "_id": 1}
	classes, err := r.CollRepo.GetWithProjection(ctx, filter, projection)
	if err != nil {
		return []any{}, fmt.Errorf("failed to get classes by email: %w", err)
	}

	return classes, nil
}

// UpdateClass implements repository.ClassRepository.UpdateClass
func (r *ClassMongoRepository) UpdateClass(ctx context.Context, class *entity.Class) (any, error) {
	filter := bson.M{"_id": class.ID} // Giả sử class có trường ID kiểu ObjectID
	fmt.Println("CLASS.StudentsWait: ", class.StudentsWait)
	_, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": class})
	if err != nil {
		return nil, fmt.Errorf("failed to update class: %w", err)
	}
	return class, nil
}

// DeleteClass implements repository.ClassRepository.DeleteClass
func (r *ClassMongoRepository) DeleteClass(ctx context.Context, emailID string, id primitive.ObjectID) error {
	filter := bson.M{"email_id": emailID, "_id": id}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}
	return nil
}

func (r *ClassMongoRepository) JoinClass(ctx context.Context, classID primitive.ObjectID, email string) error {
	// Tạo filter để tìm class theo _id
	filter := bson.M{"_id": classID}

	// Tìm class dựa trên classID để kiểm tra giá trị của is_public
	result, err := r.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find class: %w", err)
	}

	if result == nil {
		return fmt.Errorf("no class found with the given ID")
	}

	// Thực hiện type assertion để đảm bảo result có kiểu map[string]interface{}
	classDoc, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for class document, expected map[string]interface{} but got %T", result)
	}

	// Kiểm tra giá trị của is_public
	isPublic, ok := classDoc["is_public"].(bool)
	if !ok {
		return fmt.Errorf("is_public field not found or is not of type bool, found type: %T", classDoc["is_public"])
	}

	// Khởi tạo students_accept và students_wait nếu chúng là null
	update := bson.M{}
	if studentsAccept, ok := classDoc["students_accept"]; !ok || studentsAccept == nil {
		update["$set"] = bson.M{"students_accept": bson.A{}}
	}
	if studentsWait, ok := classDoc["students_wait"]; !ok || studentsWait == nil {
		update["$set"] = bson.M{"students_wait": bson.A{}}
	}

	// Thực hiện cập nhật ban đầu nếu cần
	if len(update) > 0 {
		_, err = r.CollRepo.Update(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to initialize students fields: %w", err)
		}
	}

	// Cập nhật students_accept hoặc students_wait dựa trên is_public
	if isPublic {
		_, err = r.CollRepo.Update(ctx, filter, bson.M{
			"$addToSet": bson.M{"students_accept": email},
		})
	} else {
		_, err = r.CollRepo.Update(ctx, filter, bson.M{
			"$addToSet": bson.M{"students_wait": email},
		})
	}

	if err != nil {
		return fmt.Errorf("failed to update class: %w", err)
	}

	return nil
}

func (r *ClassMongoRepository) GetQuestionOfTest(ctx context.Context, classID, testID primitive.ObjectID, email string) ([]primitive.ObjectID, primitive.M, error) {
	filter := bson.M{
		"_id": classID,
		"test": bson.M{
			"$elemMatch": bson.M{
				"_id": testID,
			},
		},
	}

	projection := bson.M{
		"test.$": 1, // Lấy đúng test cần tìm trong mảng
	}

	result, err := r.CollRepo.GetOneWithProjection(ctx, filter, projection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get test from class: %w", err)
	}
	fmt.Printf("Result: %#v\n", result) // In kết quả trả về để kiểm tra

	if len(result) == 0 {
		return nil, nil, fmt.Errorf("no class found")
	}

	tests, ok := result["test"].(bson.A)
	if !ok || len(tests) == 0 {
		return nil, nil, fmt.Errorf("no matching test found")
	}

	testDoc, ok := tests[0].(bson.M)
	if !ok {
		return nil, nil, fmt.Errorf("invalid test document format")
	}

	questionIDsRaw, ok := testDoc["question_ids"].(bson.A)
	if !ok {
		return nil, nil, fmt.Errorf("question_ids not found or invalid")
	}
	fmt.Println(questionIDsRaw)

	var questionIDs []primitive.ObjectID
	for _, q := range questionIDsRaw {
		if oid, ok := q.(primitive.ObjectID); ok {
			questionIDs = append(questionIDs, oid)
		}
	}

	delete(testDoc, "question_ids")

	return questionIDs, testDoc, nil
}

func (r *ClassMongoRepository) GetAllTestOfClass(ctx context.Context, email string, ids primitive.ObjectID) ([]any, error) {
	filter := bson.M{
		"_id": ids,
	}

	projection := bson.M{
		"test":       1,
		"class_name": 1,
		"_id":        1,
	}

	classes, err := r.CollRepo.GetWithProjection(ctx, filter, projection)
	if err != nil {
		return nil, fmt.Errorf("failed to get class tests: %w", err)
	}

	var tests []any

	for _, c := range classes {
		classDoc, ok := c.(bson.M)
		if !ok {
			continue
		}

		className := classDoc["class_name"]
		classID := classDoc["_id"]
		classTests, ok := classDoc["test"].(bson.A)
		if !ok {
			continue
		}

		for _, t := range classTests {
			testDoc, ok := t.(bson.M)
			if !ok {
				continue
			}

			testDoc["class_name"] = className
			testDoc["class_id"] = classID

			tests = append(tests, testDoc) // testDoc là kiểu bson.M → phù hợp với any
		}
	}

	return tests, nil

}
