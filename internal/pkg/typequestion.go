package pkg

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TypeQuestion trả về tên trường cho danh sách và câu trả lời dựa trên loại câu hỏi
func TypeQuestion(typeQuestion string) (string, string) {
	switch typeQuestion {
	case "fill_in_the_blank":
		return "fill_in_the_blank", "correct_answer"

	case "single_choice_question", "multiple_choice_question":
		return "options", "iscorrect"

	case "order_question":
		return "options", "order"
	case "match_choice_question":
		return "options", "match"
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

func ProcessQuestion(question bson.M) bson.M {
	// Determine the field names for options and answers based on the question type
	typeQuestion := question["type"].(string)
	arrayField, answerField := TypeQuestion(typeQuestion)
	switch typeQuestion {
	case "match_choice_question":
		if array, ok := question[arrayField].(bson.A); ok {
			//Swap random array match field
			swapFields := swapMatchFields(array)
			// Update the question with modified options
			question[arrayField] = swapFields
		}
		return question
	default:
		if arrayField == "" || answerField == "" {
			return question // Return original question if fields are not determined
		}

		switch v := question[arrayField].(type) {
		case primitive.A:
			// Handle primitive.A (bson.A is an alias for []interface{})
			convertedArray := ConvertToBsonMArray(question[arrayField].(bson.A))
			convertedArray = RemoveAnswer(convertedArray, answerField)
			question[arrayField] = convertedArray

		case []interface{}:
			// Handle []interface{} directly
			convertedArray := ConvertToBsonMArray(question[arrayField].([]interface{}))
			convertedArray = RemoveAnswer(convertedArray, answerField)
			question[arrayField] = convertedArray

		case []bson.M:
			// Handle []bson.M directly
			question[arrayField] = RemoveAnswer(question[arrayField].([]bson.M), answerField)

		default:
			fmt.Printf("Unsupported array type: %T\n", v)
		}

		// Check if the field is a bson.A (array in MongoDB terms)

		return question // Return the modified question
	}
}

func swapMatchFields(options bson.A) bson.A {
	for i := 0; i < len(options)/2; i++ {
		// Calculate the index of the matching pair to swap
		j := len(options) - 1 - i
		// Type assert each option to bson.M before accessing the fields
		if option1, ok1 := options[i].(map[string]interface{}); ok1 {
			if option2, ok2 := options[j].(map[string]interface{}); ok2 {
				// Swap the match fields
				option1["match"], option2["match"] = option2["match"], option1["match"]
				options[i], options[j] = option2, option1
			}
		}

	}
	return options
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
