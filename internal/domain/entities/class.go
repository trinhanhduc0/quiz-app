package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	ClassName     string               `json:"class_name" bson:"class_name"`
	AuthorMail    string               `json:"author_mail" bson:"author_mail"`
	TestID        []primitive.ObjectID `json:"test_id" bson:"test_id"`
	EmailID       string               `json:"email_id" bson:"email_id"`
	StudentAccept []string             `json:"students_accept" bson:"students_accept"`
	StudentsWait  []string             `json:"students_wait" bson:"students_wait"`
	IsPublic      bool                 `json:"is_public" bson:"is_public"`
	UpdatedAt     time.Time            `json:"updated_at" bson:"updated_at"`
	CreatedAt     time.Time            `json:"created_at" bson:"created_at"`
	Tags          []string             `json:"tags" bson:"tags"`
	CodeClass     string
}

func (c *Class) Validate() error {
	if c.ID.IsZero() {
		return errors.New("invalid Class ID")
	}

	// Check list TestID
	for _, TestID := range c.TestID {
		if TestID.IsZero() {
			return errors.New("invalid Test ID")
		}
	}

	return nil
}

func CreateNewClass(newClass Class) (Class, error) {
	if err := newClass.Validate(); err != nil {
		return Class{}, err
	}

	newClass.UpdatedAt = time.Now()
	return newClass, nil
}
