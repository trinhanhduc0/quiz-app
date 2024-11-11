package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/pkg"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesQuestion struct {
	auth            service.AuthHandler
	questionUseCase service.QuestionUseCase
}

func NewRouterQuestion(questionUseCase service.QuestionUseCase, auth service.AuthHandler) *RoutesQuestion {
	return &RoutesQuestion{
		questionUseCase: questionUseCase,
		auth:            auth,
	}
}

func (rq *RoutesQuestion) GetQuestionRouter(r *Router) {
	r.Handle("/questions", rq.auth.AuthMiddleware(http.HandlerFunc(rq.createQuestions))).Methods("POST")
	r.Handle("/questions", rq.auth.AuthMiddleware(http.HandlerFunc(rq.getAllQuestions))).Methods("GET")
	r.Handle("/questions", rq.auth.AuthMiddleware(http.HandlerFunc(rq.updateQuestion))).Methods("PATCH")
	r.Handle("/questions", rq.auth.AuthMiddleware(http.HandlerFunc(rq.deleteQuestion))).Methods("DELETE")
}

func (r *RoutesQuestion) createQuestions(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("email_id").(string)
	var question entity.Question

	if err := json.NewDecoder(req.Body).Decode(&question); err != nil {
		pkg.SendError(w, "Question not created", http.StatusInternalServerError)
		return
	}
	// Tạo ID mới cho các câu hỏi dựa trên loại câu hỏi
	switch question.Type {
	case "fill_in_the_blank":
		for i := range question.FillInTheBlanks {
			question.FillInTheBlanks[i].ID = primitive.NewObjectID()
			fmt.Println(question.FillInTheBlanks[i].ID)
		}
	case "match_choice_question":
		for i := range question.Options {
			question.Options[i].MatchId = primitive.NewObjectID()
			question.Options[i].ID = primitive.NewObjectID()
		}
	case "multiple_choice_question", "single_choice_question", "order_question":
		for i := range question.Options {
			question.Options[i].ID = primitive.NewObjectID()
		}
	default:
		fmt.Println("Error: Unknown question type")
	}

	// Add metadata and timestamps
	question.Metadata.Author = userID
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	insertedQuestion, err := r.questionUseCase.CreateQuestion(context.TODO(), &question)

	if err != nil {
		pkg.SendError(w, "Question not created", http.StatusInternalServerError)
		return
	}
	// Send a successful response with the inserted ID
	pkg.SendResponse(w, http.StatusCreated, insertedQuestion)
}

func (rq *RoutesQuestion) getAllQuestions(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("email_id").(string)

	// Parse limit and page from query parameters, default to 50 and 0 if not provided
	limitParam := req.URL.Query().Get("limit")
	pageParam := req.URL.Query().Get("page")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 50 // Default to 50 items per page
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 0 {
		page = 0 // Default to page 0 if not specified
	}

	questions, err := rq.questionUseCase.GetAllQuestionsByUser(context.TODO(), userID, limit, page)
	if err != nil {
		pkg.SendError(w, "Failed to get questions", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, questions)
}

func (r *RoutesQuestion) updateQuestion(w http.ResponseWriter, req *http.Request) {
	emailID := req.Context().Value("email_id").(string)

	var question entity.Question
	if err := json.NewDecoder(req.Body).Decode(&question); err != nil {
		pkg.SendError(w, "Question not updated", http.StatusInternalServerError)
		return
	}

	question.Metadata.Author = emailID
	switch question.Type {
	case "fill_in_the_blank":
		// Giả sử bạn có một danh sách các câu trả lời từ người dùng
		for i, v := range question.FillInTheBlanks {
			if primitive.NilObjectID == v.ID {
				question.FillInTheBlanks[i].ID = primitive.NewObjectID()
			}
		}
	case "match_choice_question":
		for i, v := range question.Options {
			if primitive.NilObjectID == v.ID {
				question.Options[i].ID = primitive.NewObjectID()
			}
			fmt.Println(v.MatchId)
			if primitive.NilObjectID == v.MatchId {
				question.Options[i].MatchId = primitive.NewObjectID()
			}
			fmt.Println(question.Options[i].MatchId)
		}
	case "multiple_choice_question", "single_choice_question", "order_question":
		// Giả sử bạn có một danh sách các câu trả lời từ người dùng
		for i, v := range question.Options {
			if primitive.NilObjectID == v.ID {
				question.Options[i].ID = primitive.NewObjectID()
			}
		}

	default:
		fmt.Println("Error: Unknown question type")
	}

	questionUpdated, err := r.questionUseCase.UpdateQuestion(context.TODO(), &question)
	if err != nil {
		pkg.SendError(w, "Failed to update question", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, questionUpdated)

}

func (rq *RoutesQuestion) deleteQuestion(w http.ResponseWriter, req *http.Request) {
	emailID := req.Context().Value("email_id").(string)
	var question entity.Question
	if err := json.NewDecoder(req.Body).Decode(&question); err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	question.Metadata.Author = emailID
	err := rq.questionUseCase.DeleteQuestion(context.TODO(), &question)
	if err != nil {
		pkg.SendError(w, "Failed to delete question", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, question.ID)
}
