package routes

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	entity "quiz-app/internal/domain/entities"
// 	"quiz-app/internal/domain/service"
// 	"quiz-app/internal/pkg"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Question struct {
// 	ID              string       `json:"_id"`
// 	CreatedAt       string       `json:"created_at"`
// 	Options         []Option     `json:"options"`
// 	QuestionContent QuestionText `json:"question_content"`
// 	Score           float64      `json:"score"`
// 	Tags            []string     `json:"tags"`
// 	Type            string       `json:"type"`
// 	UpdatedAt       string       `json:"updated_at"`
// }

// type Metadata struct {
// 	Author string `json:"author"`
// }

// type Option struct {
// 	ID        string `json:"id"`
// 	IsCorrect bool   `json:"iscorrect"`
// 	Text      string `json:"text"`
// 	Match     string `json:"match"`
// 	ImageURL  string `json:"imageurl"`
// 	Order     int    `json:"order"`
// }

// type QuestionText struct {
// 	Text string `json:"text"`
// }

// type QuestionsFromRedis []Question

// type RouterAnswer struct {
// 	auth *service.AuthHandler

// 	testUseCase     *service.TestUseCase
// 	answerUseCase   *service.AnswerUseCase
// 	questionUseCase *service.QuestionUseCase
// 	redisUseCase    *service.RedisUseCase
// }

// func NewRouterAnswer(s *service.AnswerUseCase, t *service.TestUseCase, q *service.QuestionUseCase, r *service.RedisUseCase, auth *service.AuthHandler) RouterAnswer {
// 	return RouterAnswer{
// 		auth: auth,

// 		answerUseCase:   s,
// 		testUseCase:     t,
// 		redisUseCase:    r,
// 		questionUseCase: q,
// 	}
// }

// func (rc RouterAnswer) GetAnswerRouter(r *Router) {
// 	r.Router.Handle("/answer/update", rc.auth.AuthMiddleware(http.HandlerFunc(rc.updateAnswer))).Methods("POST")

// 	r.Router.Handle("/answer/get", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAnswer))).Methods("POST")

// 	r.Router.Handle("/answer/user", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAllAnswerByEmail))).Methods("GET")

// }

// func (rc RouterAnswer) updateAnswer(w http.ResponseWriter, req *http.Request) {
// 	emailId := req.Context().Value("email_id").(string)
// 	email := req.Context().Value("email").(string)

// 	var newAnswer entity.TestAnswer
// 	if err := json.NewDecoder(req.Body).Decode(&newAnswer); err != nil {
// 		fmt.Println(err)
// 		pkg.SendError(w, "Invalid answer field", http.StatusBadRequest)
// 		return
// 	}

// 	newAnswer.EmailID = emailId
// 	newAnswer.Email = email

// 	idQuestions, err := rc.redisUseCase.HGetAll(req.Context(), newAnswer.TestId.Hex()+"_id")
// 	idOptions, err := rc.redisUseCase.HGetAll(req.Context(), newAnswer.TestId.Hex())
// 	rc.processIdOption(&newAnswer, idOptions, idQuestions)

// 	cacheKeyTestQuestion := fmt.Sprintf("questions_id_%s", newAnswer.TestId.Hex())
// 	cacheKeyTestQuestionDefault := fmt.Sprintf("questions_default_%s", newAnswer.TestId.Hex())
// 	questionsFromRedis, err := rc.redisUseCase.HGetAll(req.Context(), cacheKeyTestQuestion)
// 	if err != nil {
// 		return
// 	}
// 	questionsDefaultFromRedis, err := rc.redisUseCase.Get(req.Context(), cacheKeyTestQuestionDefault)
// 	if err != nil {
// 		return
// 	}

// 	// Assuming your questions are stored under a specific key like "questions"
// 	questionsJson, ok := questionsFromRedis["questions"]
// 	if !ok {
// 		fmt.Println("No questions found in Redis")
// 		return
// 	}

// 	fmt.Println("questionsDefaultFromRedis: ", questionsDefaultFromRedis)

// 	questions, err := decodeQuestions([]byte(questionsJson))
// 	if err != nil {
// 		fmt.Println("Failed to decode questions:", err)
// 		return
// 	}

// 	newAnswer.TotalScore = float32(rc.processScore(questions[newAnswer.TestId.Hex()], &newAnswer))

// 	err = rc.answerUseCase.UpdateAnswer(req.Context(), newAnswer)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to create answer: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	pkg.SendResponse(w, http.StatusCreated, newAnswer)
// }

// // decodeQuestions decodes a JSON byte array into a map[string][]map[string]interface{}
// func decodeQuestions(questionsJson []byte) (map[string]map[string]map[string]interface{}, error) {
// 	var questionsMap map[string]map[string]map[string]interface{}
// 	err := json.Unmarshal(questionsJson, &questionsMap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return questionsMap, nil
// }

// func (rc RouterAnswer) processIdOption(answer *entity.TestAnswer, IdOption map[string]string, IdQuestion map[string]string) {
// 	// Iterate through the list of questions in the answer
// 	for i, question := range answer.ListQuestionAnswer {
// 		// Convert the question ObjectID to its hex string
// 		questionID := question.QuestionID.Hex()
// 		// Check if the question ID exists in the IdOption map (from Redis)
// 		if newQuestionID, ok := IdQuestion[questionID]; ok {
// 			answer.ListQuestionAnswer[i].QuestionID, _ = primitive.ObjectIDFromHex(newQuestionID)
// 			// Check if the question has options and update the option IDs
// 			if len(question.Options) > 0 || len(question.FillInTheBlanks) > 0 {
// 				// Get the option mapping from the IdOption map (which holds a JSON object)
// 				if optionMappingStr, ok := IdOption[questionID]; ok {
// 					// Unmarshal the option mapping (it's stored as a string in Redis, e.g., a JSON object)
// 					optionMapping := make(map[string]string)
// 					err := json.Unmarshal([]byte(optionMappingStr), &optionMapping)
// 					if err != nil {
// 						fmt.Println("Failed to unmarshal option mapping:", err)
// 						continue
// 					}
// 					switch answer.ListQuestionAnswer[i].Type {
// 					case "fill_in_the_blank":
// 						for j, fillInData := range question.FillInTheBlanks {
// 							oldOptionID := fillInData.ID.Hex()
// 							// If the old option ID exists in the option mapping, update the option ID
// 							if newOptionID, ok := optionMapping[oldOptionID]; ok {
// 								// Convert the new Option ID to primitive.ObjectID
// 								if newID, err := primitive.ObjectIDFromHex(newOptionID); err == nil {
// 									answer.ListQuestionAnswer[i].FillInTheBlanks[j].ID = newID
// 								} else {
// 									fmt.Printf("Error converting new Option ID to ObjectID: %v\n", err)
// 								}
// 							}
// 						}
// 					default:
// 						// Iterate through the options and replace their IDs based on the optionMapping
// 						for j, opt := range question.Options {
// 							oldOptionID := opt.ID.Hex()
// 							// If the old option ID exists in the option mapping, update the option ID
// 							if newOptionID, ok := optionMapping[oldOptionID]; ok {
// 								answer.ListQuestionAnswer[i].Options[j].ID, _ = primitive.ObjectIDFromHex(newOptionID)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (rc RouterAnswer) processScore(questions map[string]map[string]interface{}, answer *entity.TestAnswer) float64 {
// 	var totalScore float64 = 0.0

// 	// Check if there are any questions
// 	if len(questions) == 0 {
// 		fmt.Println("Error: No questions found")
// 		return totalScore
// 	}
// 	// Iterate through the user's answers
// 	for _, userAnswer := range answer.ListQuestionAnswer {
// 		// Get question data based on user's answer
// 		questionData, exists := questions[userAnswer.QuestionID.Hex()]
// 		if !exists {
// 			fmt.Println("Error: Question not found for ID:", userAnswer.QuestionID.Hex())
// 			continue
// 		}

// 		// Print and validate question type
// 		questionType, ok := questionData["type"].(string)
// 		if !ok {
// 			fmt.Println("Error: type is not of type string")
// 			continue
// 		}
// 		// Print and validate score
// 		score, ok := questionData["score"].(float64)
// 		if !ok {
// 			fmt.Println("Error: score is not of type float32")
// 			continue
// 		}
// 		// Process the question based on its type
// 		switch questionType {
// 		case "order_question":
// 			// Access the specific options for the order question
// 			if optionMap, ok := questionData["optionMap"].(map[string]interface{}); ok {
// 				fmt.Println(optionMap)
// 				fmt.Printf("%T", optionMap[userAnswer.QuestionID.Hex()])

// 				if options, ok := optionMap[userAnswer.QuestionID.Hex()].([]interface{}); ok {
// 					totalScore += rc.processOrderQuestion(userAnswer.Options, options, score)
// 				} else {
// 					fmt.Println("Error: Options not found for order question.")
// 				}
// 				fmt.Println(totalScore)
// 			} else {
// 				fmt.Println("Error: optionMap is not of type map[string]interface{}")
// 			}
// 		case "multiple_choice_question":
// 			// Access the specific options for the order question

// 			if optionMap, ok := questionData["optionMap"].(map[string]interface{}); ok {
// 				fmt.Println(optionMap)
// 				fmt.Printf("%T", optionMap[userAnswer.QuestionID.Hex()])
// 				if options, ok := optionMap[userAnswer.QuestionID.Hex()].([]interface{}); ok {
// 					totalScore += rc.processMultipleChoice(userAnswer.Options, options, score)
// 				} else {
// 					fmt.Println("Error: Options not found for multi question.")
// 				}
// 				fmt.Println(totalScore)

// 			} else {
// 				fmt.Println("Error: optionMap is not of type map[string]interface{}")
// 			}
// 		case "fill_in_the_blank":
// 			// Access the specific options for the order question
// 			if optionMap, ok := questionData["optionMap"].(map[string]interface{}); ok {
// 				if options, ok := optionMap[userAnswer.QuestionID.Hex()].([]interface{}); ok {
// 					totalScore += rc.processFillInTheBlank(userAnswer.FillInTheBlanks, options, score)
// 				} else {
// 					fmt.Println("Error: Options not found for fill in blank question.")
// 				}
// 				fmt.Println(totalScore)

// 			} else {
// 				fmt.Println("Error: optionMap is not of type map[string]interface{}")
// 			}
// 		case "match_choice_question":
// 			// Access the specific options for the order question
// 			if optionMap, ok := questionData["optionMap"].(map[string]interface{}); ok {
// 				if options, ok := optionMap[userAnswer.QuestionID.Hex()].(map[string]interface{}); ok {
// 					totalScore += rc.processMatchChoice(userAnswer.Options, options, score)
// 				} else {
// 					fmt.Println("Error: Options not found for match question.")
// 				}
// 				fmt.Println(totalScore)

// 			} else {
// 				fmt.Println("Error: optionMap is not of type map[string]interface{}")
// 			}
// 		case "single_choice_question":
// 			if optionMap, ok := questionData["optionMap"].(map[string]interface{}); ok {
// 				if options, ok := optionMap[userAnswer.QuestionID.Hex()].([]interface{}); ok {
// 					totalScore += rc.processSingleChoice(userAnswer.Options, options, score)
// 				} else {
// 					fmt.Println("Error: Options not found for single question.")
// 				}
// 				fmt.Println(totalScore)

// 			} else {
// 				fmt.Println("Error: optionMap is not of type map[string]interface{}")
// 			}
// 		default:
// 			fmt.Println("Unknown question type:", questionType)
// 		}

// 	}
// 	return totalScore
// }

// func (rc RouterAnswer) processOrderQuestion(userOptions []entity.OptionAnswer, questionData []interface{}, score float64) float64 {
// 	for i, v := range questionData {
// 		if v != userOptions[i].ID.Hex() {
// 			return 0
// 		}
// 	}
// 	return score
// }

// func (rc RouterAnswer) processMultipleChoice(userOptions []entity.OptionAnswer, questionData []interface{}, score float64) float64 {
// 	fmt.Println("Option: ", userOptions)
// 	fmt.Println("questionData: ", questionData)

// 	// Tạo một bộ để lưu trữ các ID của tùy chọn đã chọn
// 	selectedOptionIDs := make(map[string]struct{})
// 	for _, option := range userOptions {
// 		selectedOptionIDs[option.ID.Hex()] = struct{}{}
// 	}

// 	matchingCount := 0

// 	// Kiểm tra các questionData và đếm số lượng khớp
// 	for _, optionID := range questionData {
// 		optionId, err := primitive.ObjectIDFromHex(optionID.(string))

// 		fmt.Println(optionID)
// 		fmt.Println(optionId)

// 		if err == nil {
// 			if _, exists := selectedOptionIDs[optionId.Hex()]; exists {
// 				matchingCount++
// 			}
// 		}

// 	}

// 	// Kiểm tra số lượng khớp có bằng số lượng tùy chọn không
// 	if matchingCount != len(questionData) {
// 		fmt.Println("FALSE")
// 		return 0
// 	}

// 	return score
// }

// func (rc RouterAnswer) processFillInTheBlank(userAnswers []entity.FillInTheBlank, questionData []interface{}, score float64) float64 {
// 	totalScore := float64(0)

// 	// Tạo một bản đồ để dễ dàng tra cứu câu trả lời đúng theo ID
// 	correctAnswers := make(map[string]string)
// 	for _, v := range questionData {
// 		if data, ok := v.(map[string]interface{}); ok {
// 			if id, ok := data["id"].(string); ok {
// 				if correctAnswer, ok := data["correct_answer"].(string); ok {
// 					correctAnswers[id] = correctAnswer // Lưu ID và câu trả lời đúng vào bản đồ
// 				}
// 			}
// 		}
// 	}

// 	// Lặp qua các câu trả lời của người dùng
// 	for _, userAnswer := range userAnswers {
// 		if correctAnswer, exists := correctAnswers[userAnswer.ID.Hex()]; exists {
// 			fmt.Println("correctAnswer, userAnswer", correctAnswer, userAnswer.CorrectAnswer)
// 			if correctAnswer == userAnswer.CorrectAnswer {
// 				fmt.Println(correctAnswer, userAnswer)
// 				totalScore += score / float64(len(userAnswers)) // Cộng điểm nếu câu trả lời đúng
// 			}
// 		}
// 	}

// 	return totalScore
// }

// func (rc RouterAnswer) processMatchChoice(userOptions []entity.OptionAnswer, questionData map[string]interface{}, score float64) float64 {
// 	if len(userOptions) != len(questionData) {
// 		return 0
// 	}

// 	for _, userOption := range userOptions {
// 		fmt.Println(questionData[userOption.MatchId.Hex()])
// 		fmt.Println(userOption.ID.Hex())
// 		if questionData[userOption.MatchId.Hex()] != userOption.ID.Hex() {
// 			return 0
// 		}
// 	}
// 	fmt.Println(score)
// 	return score
// }

// func (rc RouterAnswer) processSingleChoice(userOptions []entity.OptionAnswer, questionData []interface{}, score float64) float64 {

// 	if userOptions[0].ID.Hex() != questionData[0] {
// 		return 0
// 	}
// 	return score
// }

// func (rc RouterAnswer) getAnswer(w http.ResponseWriter, req *http.Request) {

// }

// func (rc RouterAnswer) getAllAnswerByEmail(w http.ResponseWriter, req *http.Request) {
// 	email := req.Context().Value("email").(string)

// 	answers, err := rc.answerUseCase.GetAllAnswerByEmail(req.Context(), email)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to retrieve answers: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	pkg.SendResponse(w, http.StatusOK, answers)
// }
