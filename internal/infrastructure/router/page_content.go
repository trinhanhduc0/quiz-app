package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/pkg"
	utils "quiz-app/internal/util"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesTest struct {
	auth service.AuthHandler

	redisUseCase service.RedisUseCase

	testUseCase     service.TestUseCase
	questionUseCase service.QuestionUseCase
	classUseCase    service.ClassUseCase
	answerUseCase   service.AnswerUseCase
}

// NewRouterTest creates a new RoutesTest instance
func NewRouterTest(testUseCase service.TestUseCase, classUseCase service.ClassUseCase, questionUseCase service.QuestionUseCase, answerUseCase service.AnswerUseCase, redisUseCase service.RedisUseCase, auth service.AuthHandler) RoutesTest {
	return RoutesTest{
		testUseCase:     testUseCase,
		classUseCase:    classUseCase,
		questionUseCase: questionUseCase,
		redisUseCase:    redisUseCase,
		answerUseCase:   answerUseCase,
		auth:            auth,
	}
}

// InitializeRoutesTests initializes all test-related routes
func (rt RoutesTest) GetTestRouter(r *Router) {
	// Routes for managing tests
	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getAllTestFromAuthor))).Methods("GET")
	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.createTest))).Methods("POST")
	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.updateTest))).Methods("PATCH")
	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.deleteTest))).Methods("DELETE")

	// Routes for class-specific operations
	r.Handle("/tests/class", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getAllTestOfClassByEmail))).Methods("POST")

	// Routes for managing questions within tests
	r.Handle("/tests/questions", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getQuestionOfTest))).Methods("POST")

	// Routes for marking test completion and sending test results
	r.Handle("/test/done", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getDoneTest))).Methods("POST")
}

func (r *RoutesTest) createTest(w http.ResponseWriter, req *http.Request) {
	emailID := req.Context().Value("email_id").(string)
	email := req.Context().Value("email").(string)

	var test entity.Test

	// Generate update fields from the test struct
	if err := json.NewDecoder(req.Body).Decode(&test); err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	test.EmailID = emailID
	test.EmailName = email

	insertedID, err := r.testUseCase.CreateTest(context.TODO(), &test)

	if err != nil {
		pkg.SendError(w, "Invalid create test", http.StatusBadRequest)
	}

	test.ID = insertedID
	pkg.SendResponse(w, http.StatusCreated, test)
}

func (r *RoutesTest) updateTest(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value("email_id").(string)

	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	var testUpdate entity.Test
	// Generate update fields from the test struct
	if err := json.NewDecoder(req.Body).Decode(&testUpdate); err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	testUpdate.EmailID = emailID
	r.testUseCase.UpdateTest(context.TODO(), &testUpdate)
	pkg.SendResponse(w, http.StatusOK, testUpdate)
}

func (r *RoutesTest) deleteTest(w http.ResponseWriter, req *http.Request) {
	emailID := req.Context().Value("email_id").(string)

	var testDelete struct {
		ID primitive.ObjectID `json:"_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&testDelete); err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := r.testUseCase.DeleteTest(context.TODO(), testDelete.ID, emailID)
	if err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	pkg.SendResponse(w, http.StatusOK, "")
}

// GetAllTestOfClassByEmail retrieves all tests for a specific class by email and class ID.
func (r *RoutesTest) getAllTestFromAuthor(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value("email").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	// Fetch tests based on class ID and email
	tests, err := r.testUseCase.GetTestsByAuthorEmail(req.Context(), emailID)
	if err != nil {
		pkg.SendError(w, "Failed to get tests", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, tests)
}

// GetAllTestOfClassByEmail retrieves all tests for a specific class by email and class ID.
func (r *RoutesTest) getAllTestOfClassByEmail(w http.ResponseWriter, req *http.Request) {
	// Extract email from context
	email, ok := req.Context().Value("email").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	type classIDRequest struct {
		ClassIDs primitive.ObjectID `json:"_id"`
	}

	// Decode class IDs from request body
	var classIDData classIDRequest
	if err := json.NewDecoder(req.Body).Decode(&classIDData); err != nil {
		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Validate that ClassIDs is not empty
	if len(classIDData.ClassIDs) == 0 {
		pkg.SendError(w, "Class ID list cannot be empty", http.StatusBadRequest)
		return
	}

	// Fetch tests based on class IDs and email
	tests, err := r.classUseCase.GetAllTestOfClass(req.Context(), email, classIDData.ClassIDs)
	if err != nil {
		pkg.SendError(w, "Failed to get tests", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, tests)
}

func (r *RoutesTest) getQuestionOfTest(w http.ResponseWriter, req *http.Request) {
	email, ok := req.Context().Value("email").(string)
	emailID, ok := req.Context().Value("email_id").(string)

	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	var test struct {
		ClassID   primitive.ObjectID `json:"class_id"`
		TestID    primitive.ObjectID `json:"test_id"`
		EmailName string             `json:"author_mail"`
		IsTest    bool               `json:"is_test"`
	}

	if err := json.NewDecoder(req.Body).Decode(&test); err != nil {
		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	cacheKeyTestInfo := fmt.Sprintf("test_info_%s", test.TestID.Hex())
	cacheKeyQuestions := fmt.Sprintf("questions_%s", test.TestID.Hex())

	// Check and load cached test info and questions
	testInfo, questions, err := r.loadCachedTestData(req.Context(), cacheKeyTestInfo, cacheKeyQuestions, email)
	if err == nil && testInfo != nil && questions != nil {
		pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": pkg.ShuffleQuestionsAndAnswers(questions)})
		return
	}

	// Fetch question IDs and test info if not cached
	questionIDs, testInfo, err := r.classUseCase.GetQuestionOfTest(req.Context(), test.ClassID, test.TestID, email)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve test data", http.StatusInternalServerError)
		return
	}
	// Cache test info

	var minutesDuration int
	if val, ok := testInfo["duration_minutes"].(int64); ok {
		minutesDuration = int(val)
	} else {
		fmt.Println("Fail to get testInfo['duration_minutes']")
	}
	// Convert minutesDuration directly to time.Duration
	duration := time.Duration(minutesDuration+1) * time.Minute
	cacheData(r.redisUseCase, req.Context(), cacheKeyTestInfo, testInfo, duration)
	// Delete allowed users for security and validate timing
	if !isTestAccessible(testInfo) {
		pkg.SendError(w, "TEST IS NOT ALLOWED", http.StatusForbidden)
		return
	}

	// Fetch questions and cache results
	questions, err = r.questionUseCase.GetAllTestQuestions(req.Context(), questionIDs)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve questions", http.StatusInternalServerError)
		return
	}

	// Generate final options and cache them
	finalOptionMap := r.generateFinalOptionsMap(questions, test.TestID.Hex())
	questionsJSON, err := json.Marshal(finalOptionMap)
	if err != nil {
		pkg.SendError(w, "Failed to process questions", http.StatusInternalServerError)
		return
	}
	r.redisUseCase.HSet(req.Context(), fmt.Sprintf("questions_id_%s", test.TestID.Hex()), duration, map[string]interface{}{"questions": questionsJSON})

	if testInfo["is_test"] == true {
		shuffledQuestions := pkg.ShuffleQuestionsAndAnswers(questions)
		cacheData(r.redisUseCase, req.Context(), cacheKeyQuestions, shuffledQuestions, duration)
		if answer, err := r.answerUseCase.GetAnswer(req.Context(), primitive.M{"test_id": test.TestID, "email_id": emailID}); err == nil && len(answer.ListQuestionAnswer) != 0 {
			response := primitive.M{"test_info": testInfo, "answer": answer, "questions": questions}
			pkg.SendResponse(w, http.StatusOK, response)
			return
		} else {
			r.answerUseCase.CreateNewAnswer(req.Context(), &entity.TestAnswer{
				TestId:  test.TestID,
				EmailID: emailID,
				Email:   email,
			})
			pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": shuffledQuestions})
			return
		}
	} else {
		if answer, err := r.answerUseCase.GetAnswer(req.Context(), primitive.M{"test_id": test.TestID, "email_id": emailID}); err == nil && len(answer.ListQuestionAnswer) != 0 {
			response := primitive.M{"test_info": testInfo, "answer": answer, "questions": questions}
			pkg.SendResponse(w, http.StatusOK, response)
			return
		}
		pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": questions})
		return
	}
}

// loadCachedTestData handles the retrieval and decoding of cached data
func (r *RoutesTest) loadCachedTestData(ctx context.Context, testInfoKey, questionsKey, email string) (map[string]interface{}, []primitive.M, error) {
	cachedTestInfo, errInfo := r.redisUseCase.Get(ctx, testInfoKey)
	cachedQuestions, errQuestions := r.redisUseCase.Get(ctx, questionsKey)
	if errInfo == nil && errQuestions == nil {
		var testInfo map[string]interface{}
		var questions []primitive.M
		if err := json.Unmarshal([]byte(cachedTestInfo), &testInfo); err != nil {
			return nil, nil, err
		}
		if err := json.Unmarshal([]byte(cachedQuestions), &questions); err != nil {
			return nil, nil, err
		}
		for _, user := range testInfo["allowed_users"].([]interface{}) {
			if user == email {
				delete(testInfo, "allowed_users")
				return testInfo, questions, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("cache miss")
}

// generateFinalOptionsMap prepares the options map for questions
func (r *RoutesTest) generateFinalOptionsMap(questions []primitive.M, testID string) map[string]map[string]map[string]interface{} {
	finalOptionMap := make(map[string]map[string]map[string]interface{})
	for _, question := range questions {
		questionType, ok := question["type"].(string)
		if !ok {
			continue
		}
		questionID := r.getQuestionID(question)
		if questionID == "" {
			continue
		}

		optionMap := map[string]interface{}{questionID: []string{}}
		switch questionType {
		case "order_question":
			optionMap[questionID] = r.handleOrderQuestion(question)
		case "single_choice_question":
			optionMap[questionID] = r.handleSingleChoiceQuestion(question)
		case "multiple_choice_question":
			optionMap[questionID] = r.handleMultipleChoiceQuestion(question)
		case "fill_in_the_blank":
			optionMap[questionID] = r.handleFillInTheBlank(question)
		case "match_choice_question":
			optionMap[questionID] = r.handleMatchChoiceQuestion(question)
		}

		if finalOptionMap[testID] == nil {
			finalOptionMap[testID] = make(map[string]map[string]interface{})
		}
		finalOptionMap[testID][questionID] = map[string]interface{}{
			"optionMap": optionMap,
			"type":      questionType,
			"score":     question["score"].(float64),
		}
	}
	return finalOptionMap
}

// Helper to retrieve question ID
func (r *RoutesTest) getQuestionID(question primitive.M) string {
	if id, ok := question["_id"].(primitive.ObjectID); ok {
		return id.Hex()
	} else if idStr, ok := question["_id"].(string); ok {
		return idStr
	}
	return ""
}

func cacheData(redisUseCase service.RedisUseCase, ctx context.Context, key string, data interface{}, time time.Duration) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data for caching:", err)
		return
	}
	if err := redisUseCase.Set(ctx, key, dataJSON, time); err != nil {
		fmt.Println("Error caching data:", err)
	}
}

func isTestAccessible(testInfo map[string]interface{}) bool {
	startTime, errStart := utils.StringToTime(testInfo["start_time"].(string))
	endTime, errEnd := utils.StringToTime(testInfo["end_time"].(string))

	if errStart != nil || errEnd != nil {
		fmt.Println("Error processing test timing")
		return false
	}
	currentTime := time.Now()
	return startTime.Before(currentTime) && endTime.After(currentTime)
}

// // sendTest processes the test submission
// func (r *RoutesTest) sendTest(w http.ResponseWriter, req *http.Request) {

// 	// TODO: Process the test submission based on email and emailID

// 	pkg.SendResponse(w, http.StatusOK, "Test sent successfully")
// }

// getDoneTest handles retrieving a user's done test
func (r *RoutesTest) getDoneTest(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value("email_id").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	// TODO: Implement fetching done test logic

	pkg.SendResponse(w, http.StatusOK, fmt.Sprintf("Done test retrieved for email ID: %s", emailID))
}

// Handle question
// handleOrderQuestion processes "order_question" type questions
func (r *RoutesTest) handleOrderQuestion(question primitive.M) []string {
	optionsValue, exists := question["options"]
	if !exists {
		fmt.Println("Error: Options key does not exist in the question map")
		return nil
	}

	rawOptions, ok := optionsValue.(primitive.A)
	if !ok {
		fmt.Println("Error: Options is not of type primitive.A")
		return nil
	}

	options := make([]map[string]interface{}, 0, len(rawOptions))
	for _, opt := range rawOptions {
		if option, ok := opt.(primitive.M); ok {
			options = append(options, option)
		} else {
			fmt.Println("Error: Option is not of type map[string]interface{}")
		}
	}

	// Sort options based on the "order" field
	sort.Slice(options, func(i, j int) bool {
		return options[i]["order"].(int32) < options[j]["order"].(int32)
	})

	var orderedIDs []string
	for _, option := range options {
		if id, ok := option["id"].(primitive.ObjectID); ok {
			orderedIDs = append(orderedIDs, id.Hex())
		} else {
			fmt.Println("Error: Option ID is not of type ObjectID")
		}
	}
	return orderedIDs
}

// handleSingleChoiceQuestion processes "single_choice_question" type questions
func (r *RoutesTest) handleSingleChoiceQuestion(question primitive.M) []string {
	optionsValue, ok := question["options"].(primitive.A)
	if !ok {
		fmt.Println("Error: Options is not of type primitive.A")
		return nil
	}

	for _, option := range optionsValue {
		if optionMap, ok := option.(primitive.M); ok {

			if isCorrect, exists := optionMap["iscorrect"].(bool); exists && isCorrect {
				if id, idOk := optionMap["id"].(primitive.ObjectID); idOk {
					return []string{id.Hex()}
				}
			}
		}
	}
	return nil
}

// handleMultipleChoiceQuestion processes "multiple_choice_question" type questions
func (r *RoutesTest) handleMultipleChoiceQuestion(question primitive.M) []string {
	optionsValue, ok := question["options"].(primitive.A)
	if !ok {
		fmt.Println("Error: Options is not of type primitive.A")
		return nil
	}

	var correctIDs []string
	for _, option := range optionsValue {
		if optionMap, ok := option.(primitive.M); ok {
			if isCorrect, exists := optionMap["iscorrect"].(bool); exists && isCorrect {
				if id, idOk := optionMap["id"].(primitive.ObjectID); idOk {
					correctIDs = append(correctIDs, id.Hex())
				}
			}
		}
	}
	return correctIDs
}

// handleFillInTheBlank processes "fill_in_the_blank" type questions
func (r *RoutesTest) handleFillInTheBlank(question primitive.M) []map[string]string {
	fillInBlanks, ok := question["fill_in_the_blank"].(primitive.A)
	if !ok {
		fmt.Println("Error: fill_in_the_blank is not of type primitive.A")
		return nil
	}

	var fillInData []map[string]string
	for _, item := range fillInBlanks {
		itemMap, ok := item.(primitive.M)
		if !ok {
			fmt.Println("Error: item in fill_in_the_blank is not of type map[string]interface{}")
			continue
		}

		id, idOk := itemMap["id"].(primitive.ObjectID)
		answer, answerOk := itemMap["correct_answer"].(string)
		if idOk && answerOk {
			fillInData = append(fillInData, map[string]string{
				"id":             id.Hex(),
				"correct_answer": answer,
			})
		}
	}
	return fillInData
}

// handleMatchChoiceQuestion processes "match_choice_question" type questions
func (r *RoutesTest) handleMatchChoiceQuestion(question primitive.M) map[string]string {
	options, ok := question["options"].(primitive.A)
	if !ok {
		fmt.Println("Error: Options is not of type primitive.A")
		return nil
	}

	matchMap := make(map[string]string)
	for _, option := range options {
		optionMap, ok := option.(primitive.M)
		if !ok {
			fmt.Println("Error: Option in match_choice_question is not of type map[string]interface{}")
			continue
		}

		matchID, matchIDOk := optionMap["matchid"].(primitive.ObjectID)
		id, idOk := optionMap["id"].(primitive.ObjectID)
		if matchIDOk && idOk {
			matchMap[matchID.Hex()] = id.Hex()
		}
	}
	return matchMap
}
