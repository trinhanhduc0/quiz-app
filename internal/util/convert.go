package utils

import (
	"errors"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ArrayStringToObjectID converts a BSON array of string IDs to a slice of ObjectIDs.
func ArrayStringToObjectId(arr []any) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range arr {
		strID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", id)
		}
		objectID, err := primitive.ObjectIDFromHex(strID)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID format: %v", err)
		}
		objectIDs = append(objectIDs, objectID)
	}
	return objectIDs, nil
}

func StringToObjectId(id any) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return primitive.NilObjectID, err // Return an error instead of panicking
	}
	return objectID, nil
}

func StringToTime(timeStr string) (time.Time, error) {
	// ISO 8601 format: "2025-04-25T18:00:00.000Z"
	layout := time.RFC3339 // chuẩn ISO 8601
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid time format: " + err.Error())
	}
	return t, nil
}

// GenerateUpdateFields generates the fields to update for MongoDB based on the non-empty fields in the provided struct
func GenerateUpdateFields(targetStruct any) (bson.M, error) {
	v := reflect.ValueOf(targetStruct)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("targetStruct must be a pointer")
	}

	v = v.Elem() // Dereference the pointer to get the underlying struct
	typeOfStruct := v.Type()
	updateFields := bson.M{}

	// Iterate over the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfStruct.Field(i).Tag.Get("json")
		if fieldName == "" {
			fieldName = typeOfStruct.Field(i).Name
		}

		// Check if the field has a non-empty value
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				updateFields[fieldName] = field.String()
			}
		case reflect.Slice:
			if field.Len() > 0 {
				switch fieldName {
				case "options":
					options := cleanOptions(field.Interface().([]entity.Option))
					if len(options) > 0 {
						updateFields[fieldName] = options
					}
				case "fill_in_blanks":
					blanks := cleanFillInBlanks(field.Interface().([]entity.FillInTheBlank))
					if len(blanks) > 0 {
						updateFields[fieldName] = blanks
					}
				case "order_items":
					items := cleanOrderItems(field.Interface().([]entity.OrderItem))
					if len(items) > 0 {
						updateFields[fieldName] = items
					}
				case "match_items":
					matchItems := cleanMatchItems(field.Interface().([]entity.MatchItem))
					if len(matchItems) > 0 {
						updateFields[fieldName] = matchItems
					}
				case "match_options":
					matchOptions := cleanMatchOptions(field.Interface().([]entity.MatchOption))
					if len(matchOptions) > 0 {
						updateFields[fieldName] = matchOptions
					}
				case "correct_map":
					if !field.IsNil() && len(field.Interface().(map[string]string)) > 0 {
						updateFields[fieldName] = field.Interface()
					}
				}
			}
		case reflect.Ptr:
			if !field.IsNil() {
				updateFields[fieldName] = field.Elem().Interface()
			}
		case reflect.Struct:
			fmt.Printf("Field: %s, Value: %v\n", fieldName, field.Interface())

			// Handle nested structs like question_content
			nestedFields, err := generateNestedFields(field)
			if err != nil {
				return nil, err
			}
			if len(nestedFields) > 0 {
				updateFields[fieldName] = nestedFields
			}
		case reflect.Int:
			updateFields[fieldName] = field.Int()
		case reflect.Float32:
			updateFields[fieldName] = field.Float()
		case reflect.Bool:
			updateFields[fieldName] = field.Bool()
		}
	}

	// Return error if no fields are present
	if len(updateFields) == 0 {
		return nil, errors.New("no fields to update")
	}

	return updateFields, nil
}

// cleanOptions filters out unnecessary fields from the options slice
func cleanOptions(options []entity.Option) []bson.M {
	cleanedOptions := []bson.M{}
	for _, option := range options {
		cleanedOption := bson.M{
			"id":        option.ID,
			"text":      option.Text,
			"imageurl":  option.ImageURL,
			"iscorrect": option.IsCorrect,
		}
		cleanedOptions = append(cleanedOptions, cleanedOption)
	}
	return cleanedOptions
}

func cleanFillInBlanks(blanks []entity.FillInTheBlank) []bson.M {
	cleaned := []bson.M{}
	for _, item := range blanks {
		cleaned = append(cleaned, bson.M{
			"id":             item.ID,
			"text_before":    item.TextBefore,
			"blank":          item.Blank,
			"correct_answer": item.CorrectAnswer,
			"text_after":     item.TextAfter,
		})
	}
	return cleaned
}

func cleanOrderItems(items []entity.OrderItem) []bson.M {
	cleaned := []bson.M{}
	for _, item := range items {
		cleaned = append(cleaned, bson.M{
			"id":    item.ID,
			"text":  item.Text,
			"order": item.Order,
		})
	}
	return cleaned
}

func cleanMatchItems(items []entity.MatchItem) []bson.M {
	cleaned := []bson.M{}
	for _, item := range items {
		cleaned = append(cleaned, bson.M{
			"id":   item.ID,
			"text": item.Text,
		})
	}
	return cleaned
}

func cleanMatchOptions(items []entity.MatchOption) []bson.M {
	cleaned := []bson.M{}
	for _, item := range items {
		cleaned = append(cleaned, bson.M{
			"id":   item.ID,
			"text": item.Text,
		})
	}
	return cleaned
}

// Generate fields for nested structs
func generateNestedFields(field reflect.Value) (bson.M, error) {
	nestedFields := bson.M{}
	for j := 0; j < field.NumField(); j++ {
		nestedField := field.Field(j)
		nestedFieldType := field.Type().Field(j) //
		nestedFieldName := field.Type().Field(j).Tag.Get("json")
		if nestedFieldName == "" {
			nestedFieldName = field.Type().Field(j).Name
		}
		// Skip unexported fields
		if !nestedFieldType.IsExported() { //
			continue
		}
		if !isEmpty(nestedField) {
			nestedFields[nestedFieldName] = nestedField.Interface()
		}
	}
	return nestedFields, nil
}

// Check if a value is empty
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Slice:
		return value.Len() == 0
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	case reflect.Int:
		return value.Interface() == reflect.Zero(value.Type()).Interface()
	case reflect.Bool:
		return value.Bool()
	}

	return false
}

// ConvertIDs converts string IDs to ObjectIDs for specified fields in a map.
func ConvertIDs(fields map[string]interface{}, fieldNames ...string) {
	for _, fieldName := range fieldNames {
		if ids, ok := fields[fieldName].([]any); ok && ids != nil {
			if arrayObject, err := ArrayStringToObjectId(ids); err == nil {
				fields[fieldName] = arrayObject
			} else {
				fmt.Printf("Error converting IDs: %v\n", err)
			}
		}
	}
}

// RemoveKeysFromList removes specified keys from each map in the list
func RemoveKeysFromList(list []bson.M, keysToRemove []string) {
	for i := range list {
		for _, key := range keysToRemove {
			delete(list[i], key)
		}
	}
}

// RemoveAnswer removes specific fields from elements in an array of bson.M
func RemoveAnswer(questionList []bson.M, answerField string) {
	for _, question := range questionList {
		delete(question, answerField)
	}
}

// convertToBsonMArray chuyển đổi bson.A sang []bson.M
func ConvertToBsonMArray(array bson.A) []bson.M {
	var convertedArray []bson.M
	for _, item := range array {
		if doc, ok := item.(bson.M); ok {
			convertedArray = append(convertedArray, doc)
		}
	}
	return convertedArray
}

func RemoveEmptyFillInTheBlanks(fillInTheBlanks []entity.FillInTheBlank) []entity.FillInTheBlank {
	var result []entity.FillInTheBlank
	for _, item := range fillInTheBlanks {
		if item.CorrectAnswer != "" {
			result = append(result, item)
		}
	}
	return result
}

func RemoveEmptyOptions(options []entity.OptionAnswer) []entity.OptionAnswer {
	var result []entity.OptionAnswer
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

// Function to remove empty QuestionAnswer entries
func RemoveEmptyQuestionAnswers(answers []entity.QuestionAnswer) []entity.QuestionAnswer {
	var result []entity.QuestionAnswer
	for _, answer := range answers {
		// Clean up the fields
		answer.FillInTheBlanks = RemoveEmptyFillInTheBlanks(answer.FillInTheBlanks)
		answer.Options = RemoveEmptyOptions(answer.Options)

		// Only include non-empty QuestionAnswer entries
		if !answer.QuestionID.IsZero() && len(answer.FillInTheBlanks) > 0 || len(answer.Options) > 0 {
			result = append(result, answer)
		}
	}
	return result
}
