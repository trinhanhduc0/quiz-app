package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents the main question structure.
type Question struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,,omitempty"`
	FillInTheBlanks []FillInTheBlank   `json:"fill_in_the_blank" bson:"fill_in_the_blank,omitempty"`
	Metadata        Metadata           `json:"metadata" bson:"metadata,omitempty"`
	Options         []Option           `json:"options" bson:"options,omitempty"`
	QuestionContent QuestionContent    `json:"question_content" bson:"question_content,omitempty"`
	Suggestion      string             `json:"suggestion" bson:"suggestion,omitempty"`
	Tags            []string           `json:"tags" bson:"tags,omitempty"`
	Type            string             `json:"type" bson:"type,omitempty"`
	Score           float32            `json:"score" bson:"score,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

// FillInTheBlank represents a fill-in-the-blank part of the question.
type FillInTheBlank struct {
	ID            primitive.ObjectID `json:"id" bson:"id,omitempty"`
	TextBefore    string             `json:"text_before" bson:"text_before,omitempty"`
	Blank         string             `json:"blank" bson:"blank,omitempty"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer,omitempty"`
	TextAfter     string             `json:"text_after" bson:"text_after,omitempty"`
}

// Metadata represents the metadata for the question.
type Metadata struct {
	Author string `json:"author" bson:"author,omitempty"`
}

// Option represents each option in the question.
type Option struct {
	ID        primitive.ObjectID `json:"id" bson:"id,omitempty"`
	MatchId   primitive.ObjectID `json:"matchid" bson:"matchid,omitempty"`
	Text      string             `json:"text" bson:"text,omitempty"`
	ImageURL  string             `json:"imageurl" bson:"imageurl,omitempty"`
	IsCorrect bool               `json:"iscorrect" bson:"iscorrect,omitempty"`
	Match     string             `json:"match" bson:"match,omitempty"`
	Order     int                `json:"order" bson:"order,omitempty"`
}

// QuestionContent represents the content of the question.
type QuestionContent struct {
	Text     string `json:"text" bson:"text,omitempty"`
	ImageURL string `json:"image_url" bson:"image_url,omitempty"`
	VideoURL string `json:"video_url" bson:"video_url,omitempty"`
	AudioURL string `json:"audio_url" bson:"audio_url,omitempty"`
}
