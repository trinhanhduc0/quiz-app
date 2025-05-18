package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents the main question structure.
type Question struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type            string             `json:"type" bson:"type,omitempty"` // "fill_in_the_blank", "multiple_choice_single", ...
	QuestionContent QuestionContent    `json:"question_content,omitempty" bson:"question_content,omitempty"`
	Metadata        Metadata           `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags            []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Suggestion      []string           `json:"suggestion,omitempty" bson:"suggestion,omitempty"`
	Score           float32            `json:"score,omitempty" bson:"score,omitempty"`
	Created_At      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_At      time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

	// Polymorphic fields (chỉ dùng tương ứng với Type)
	Options         []Option         `json:"options,omitempty" bson:"options,omitempty"`                       // for multiple_choice
	FillInTheBlanks []FillInTheBlank `json:"fill_in_the_blanks,omitempty" bson:"fill_in_the_blanks,omitempty"` // for fill_in_the_blank
	OrderItems      []OrderItem      `json:"order_items,omitempty" bson:"order_items,omitempty"`               // for ordering_question
	MatchItems      []MatchItem      `json:"match_items,omitempty" bson:"match_items,omitempty"`               // for match_choice
	MatchOptions    []MatchOption    `json:"match_options,omitempty" bson:"match_options,omitempty"`           // for match_choice
	// CorrectMap      map[string]string `json:"correct_map,omitempty" bson:"correct_map,omitempty"`     // e.g. for match/map-based validation
}

// Metadata represents the metadata for the question.
type Metadata struct {
	Author string `json:"author" bson:"author,omitempty"`
}
type MatchItem struct {
	ID   primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text string             `json:"text" bson:"text,omitempty"`
}

type MatchOption struct {
	ID      primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text    string             `json:"text" bson:"text,omitempty"`
	MatchId string             `json:"match_id" bson:"match_id,omitempty"`
}

// FillInTheBlank represents a fill-in-the-blank part of the question.
type FillInTheBlank struct {
	ID            primitive.ObjectID `json:"id" bson:"id,omitempty"`
	TextBefore    string             `json:"text_before" bson:"text_before,omitempty"`
	Blank         string             `json:"blank" bson:"blank,omitempty"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer,omitempty"`
	TextAfter     string             `json:"text_after" bson:"text_after,omitempty"`
}

// Option represents each option in the question.
type Option struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Text      string             `json:"text,omitempty" bson:"text,omitempty"`
	ImageURL  string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	IsCorrect bool               `json:"iscorrect,omitempty" bson:"iscorrect,omitempty"`
}

type OrderItem struct {
	ID    primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text  string             `json:"text" bson:"text,omitempty"`
	Order int                `json:"order" bson:"order,omitempty"` // Đáp án đúng theo thứ tự
}

// QuestionContent represents the content of the question.
type QuestionContent struct {
	Text     string `json:"text" bson:"text,omitempty"`
	ImageURL string `json:"image_url" bson:"image_url,omitempty"`
	VideoURL string `json:"video_url" bson:"video_url,omitempty"`
	AudioURL string `json:"audio_url" bson:"audio_url,omitempty"`
}
