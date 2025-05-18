package pkg

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TypeQuestion trả về tên trường cho danh sách và câu trả lời dựa trên loại câu hỏi
func TypeShuffleQuestion(typeQuestion string) (string, string) {
	switch typeQuestion {
	case "fill_in_the_blank":
		return "fill_in_the_blank", "correct_answer"
	case "single_choice_question", "multiple_choice_question":
		return "options", "iscorrect"
	case "order_question":
		return "order_items", "order"
	// case "match_choice_question":
	// 	return "match_options",
	default:
		return "", ""
	}
}

// RemoveAnswer removes specific fields from elements in an array of bson.M
func RemoveAnswer(questionList []bson.M, answerField string) []bson.M {
	for _, question := range questionList {
		delete(question, answerField)
	}
	return questionList
}

// *****************SHUFFLE*******************

func ProcessQuestion(question bson.M) bson.M {
	// Determine the field names for options and answers based on the question type
	typeQuestion := question["type"].(string)
	arrayField, answerField := TypeShuffleQuestion(typeQuestion)
	switch typeQuestion {
	case "match_choice_question":
		if array, ok := question[arrayField].(bson.A); ok {
			//Swap random array match field
			swapFields := swapMatchFields(array)
			// Update the question with modified options
			question[arrayField] = swapFields
		}
		return question
	case "single_choice_question", "multiple_choice_question":
		if array, ok := question[arrayField].(bson.A); ok {
			// Swap random array options
			swapFields := swapOptions(array)
			// Update the question with modified options
			question[arrayField] = swapFields
		}
		return question
	case "order_question":
		if array, ok := question[arrayField].(bson.A); ok {
			// Swap random array order items
			swapFields := swapOrderItems(array)
			// Update the question with modified order items
			question[arrayField] = swapFields
		}
		if arrayField == "" || answerField == "" {
			return question // Return original question if fields are not determined
		}
	}
	return question
}

func swapOrderItems(array bson.A) any {
	return swap(array)
}

func swapOptions(array bson.A) bson.A {
	return swap(array)
}
func swapMatchFields(array bson.A) bson.A {
	return swap(array)
}

func swap(array bson.A) bson.A {
	n := len(array)
	result := make(bson.A, n)
	for i := 0; i < n; i++ {
		result[i] = array[n-1-i]
	}
	return result
}

func ShuffleQuestionsAndAnswers(questionList []bson.M) []bson.M {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Shuffle the list of questions
	rand.Shuffle(len(questionList), func(i, j int) {
		questionList[i], questionList[j] = questionList[j], questionList[i]
	})

	// Process each question and shuffle the options within each question
	for i := range questionList {
		questionList[i] = ProcessQuestion(questionList[i]) // Use the updated ProcessQuestion
		if options, ok := questionList[i]["options"].([]bson.M); ok {
			rand.Shuffle(len(options), func(i, j int) {
				options[i], options[j] = options[j], options[i]
			})
			// Update the options in the question after shuffling
			questionList[i]["options"] = options
		}
	}
	return questionList
}

// ConvertToBsonMArray converts bson.A to []bson.M
func ConvertToBsonMArray(array bson.A) []bson.M {
	var convertedArray []bson.M
	for _, item := range array {
		switch v := item.(type) {
		case primitive.M:
			// Convert map[string]interface{} to bson.M
			convertedArray = append(convertedArray, bson.M(v))
		case map[string]interface{}:
			convertedArray = append(convertedArray, bson.M(v))
		default:
			fmt.Printf("item is of an unexpected type: %T\n", item)
		}
	}
	return convertedArray
}
