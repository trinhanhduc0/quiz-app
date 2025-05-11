package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Test struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id"`
	TestName        string               `json:"test_name" bson:"test_name"`
	Descript        string               `json:"descript" bson:"descript"`
	QuestionIDs     []primitive.ObjectID `json:"question_ids" bson:"question_ids"`
	IsTest          bool                 `json:"is_test" bson:"is_test"`
	StartTime       string               `json:"start_time" bson:"start_time"`
	EndTime         string               `json:"end_time" bson:"end_time"`
	DurationMinutes int                  `json:"duration_minutes" bson:"duration_minutes"`
	Tags            []string             `json:"tags" bson:"tags"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at" bson:"updated_at"`
	EmailID         string               `json:"email_id" bson:"email_id"`
	EmailName       string               `json:"author_mail" bson:"author_mail"`
	AnswerUser      []string             `json:"answer_user" bson:"answer_user"` // Email user done test
}

func CreateNewTest(test Test) (Test, error) {
	// Validate UserID
	if test.ID.IsZero() {
		return Test{}, errors.New("invalid UserID")
	}

	// Validate QuestionIDs elements
	for _, questionID := range test.QuestionIDs {
		if questionID.IsZero() {
			return Test{}, errors.New("invalid QuestionID")
		}
	}

	// Set CreatedAt and UpdatedAt timestamps
	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()

	// Optionally generate a new ID if not provided
	if test.ID.IsZero() {
		test.ID = primitive.NewObjectID()
	}

	return test, nil
}
