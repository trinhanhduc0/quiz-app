package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestAnswer struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TestId             primitive.ObjectID `json:"test_id,omitempty" bson:"test_id,omitempty"`
	EmailID            string             `json:"email_id,omitempty" bson:"email_id,omitempty"`
	Email              string             `json:"email,omitempty" bson:"email,omitempty"`
	ListQuestionAnswer []QuestionAnswer   `json:"question_answer,omitempty" bson:"question_answer,omitempty"`
	CreatedAt          time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	TotalScore         float32            `json:"score" bson:"score,omitempty"`
}

type QuestionAnswer struct {
	QuestionID      primitive.ObjectID `json:"question_id,omitempty" bson:"question_id,omitempty"`
	Type            string             `json:"type,omitempty" bson:"type,omitempty"`
	FillInTheBlanks []FillInTheBlank   `json:"fill_in_the_blank,omitempty" bson:"fill_in_the_blank,omitempty"`
	Options         []OptionAnswer     `json:"options,omitempty" bson:"options,omitempty"`
}

// Option represents each option in the question.
type OptionAnswer struct {
	ID      primitive.ObjectID `json:"id" bson:"id,omitempty"`
	MatchId primitive.ObjectID `json:"matchid" bson:"matchid,omitempty"`
}

func CreateNewAnswer(answer TestAnswer) (*TestAnswer, error) {
	answer.ID = primitive.NewObjectID()
	answer.CreatedAt = time.Now()

	// Validate QuestionIDs elements
	for _, QuestionIDs := range answer.ListQuestionAnswer {
		if QuestionIDs.QuestionID.IsZero() {
			return &TestAnswer{}, errors.New("invalid QuestionID")
		}
	}

	return &answer, nil
}
