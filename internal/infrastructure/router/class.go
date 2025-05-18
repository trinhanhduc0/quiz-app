package routes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/pkg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type routerClass struct {
	auth         service.AuthHandler
	classUseCase service.ClassUseCase
	redisUseCase service.RedisUseCase
}

func NewRouterClass(classUC service.ClassUseCase, redisUC service.RedisUseCase, auth service.AuthHandler) routerClass {
	return routerClass{
		auth:         auth,
		classUseCase: classUC,
		redisUseCase: redisUC,
	}
}

// ==== Request Body Structs ====
type classIDRequest struct {
	ID     string   `json:"_id"`
	Minute int      `json:"minute"`
	TestID []string `json:"test_id"`
}

type classDeleteRequest struct {
	ID primitive.ObjectID `json:"_id"`
}

type joinClassRequest struct {
	ID string `json:"_id"`
}

// ==== Helpers ====
func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		pkg.SendError(w, "Invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}

func generateRandomKey() (string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ==== Handlers ====
func (rc routerClass) createClass(w http.ResponseWriter, req *http.Request) {
	emailID, _ := req.Context().Value("email_id").(string)
	email, _ := req.Context().Value("email").(string)

	var newClass entity.Class
	if !DecodeJSONBody(w, req, &newClass) {
		return
	}

	newClass.CreatedAt = time.Now()
	newClass.UpdatedAt = time.Now()
	newClass.AuthorMail = email
	newClass.EmailID = emailID
	newClass.StudentAccept = []string{emailID}
	newClass.StudentsWait = []string{}

	if newClass.TestID == nil {
		newClass.TestID = []primitive.ObjectID{}
	}

	createdClass, err := rc.classUseCase.CreateClass(req.Context(), &newClass)
	if err != nil {
		pkg.SendError(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, createdClass)
}

func (rc routerClass) getAllClass(w http.ResponseWriter, req *http.Request) {
	emailID, _ := req.Context().Value("email").(string)
	fmt.Println(emailID)
	classes, err := rc.classUseCase.GetAllClass(req.Context(), emailID)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}
	fmt.Println("LOG GET CLASS: ", classes)

	pkg.SendResponse(w, http.StatusOK, classes)
}

func (rc routerClass) getAllClassByEmail(w http.ResponseWriter, req *http.Request) {
	email, _ := req.Context().Value("email").(string)

	classes, err := rc.classUseCase.GetAllClassByEmail(req.Context(), email)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, classes)
}

func (rc routerClass) updateClass(w http.ResponseWriter, req *http.Request) {
	emailID, _ := req.Context().Value("email_id").(string)
	email, _ := req.Context().Value("email").(string)

	var classToUpdate entity.Class
	if !DecodeJSONBody(w, req, &classToUpdate) {
		return
	}

	classToUpdate.EmailID = emailID
	classToUpdate.AuthorMail = email

	updatedClass, err := rc.classUseCase.UpdateClass(req.Context(), &classToUpdate)
	if err != nil {
		pkg.SendError(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, updatedClass)
}

func (rc routerClass) deleteClass(w http.ResponseWriter, req *http.Request) {
	emailID, _ := req.Context().Value("email_id").(string)

	var reqBody classDeleteRequest
	if !DecodeJSONBody(w, req, &reqBody) {
		return
	}

	err := rc.classUseCase.DeleteClass(req.Context(), emailID, reqBody.ID)
	if err != nil {
		pkg.SendError(w, "Error deleting class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, fmt.Sprintf("Class %v deleted", reqBody.ID.Hex()))
}

func (rc routerClass) joinClass(w http.ResponseWriter, req *http.Request) {
	email, _ := req.Context().Value("email").(string)

	var reqBody joinClassRequest
	if !DecodeJSONBody(w, req, &reqBody) {
		return
	}

	dataJSON, err := rc.redisUseCase.Get(req.Context(), reqBody.ID)
	if err != nil {
		pkg.SendError(w, "Error fetching class from Redis", http.StatusInternalServerError)
		return
	}

	var redisData struct {
		ClassID     string   `json:"class_id"`
		EmailAuthor string   `json:"email"`
		TestID      []string `json:"test_id"`
	}
	if err := json.Unmarshal([]byte(dataJSON), &redisData); err != nil {
		pkg.SendError(w, "Error decoding Redis data", http.StatusInternalServerError)
		return
	}

	classID, err := primitive.ObjectIDFromHex(redisData.ClassID)
	if err != nil {
		pkg.SendError(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	var testIDs []primitive.ObjectID
	for i, tid := range redisData.TestID {
		id, err := primitive.ObjectIDFromHex(tid)
		if err != nil {
			pkg.SendError(w, fmt.Sprintf("Invalid test ID at index %d: %s", i, tid), http.StatusBadRequest)
			return
		}
		testIDs = append(testIDs, id)
	}

	err = rc.classUseCase.JoinClass(req.Context(), classID, testIDs, redisData.EmailAuthor, email)
	if err != nil {
		pkg.SendError(w, "Error joining class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Joined class: %s", reqBody.ID),
	})
}

func (rc routerClass) createCodeClass(w http.ResponseWriter, req *http.Request) {
	email, _ := req.Context().Value("email").(string)

	var reqBody classIDRequest
	if !DecodeJSONBody(w, req, &reqBody) {
		return
	}

	key, err := generateRandomKey()
	if err != nil {
		pkg.SendError(w, "Failed to generate code", http.StatusInternalServerError)
		return
	}

	data := struct {
		ClassID string   `json:"class_id"`
		Email   string   `json:"email"`
		TestID  []string `json:"test_id"`
	}{
		ClassID: reqBody.ID,
		Email:   email,
		TestID:  reqBody.TestID,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		pkg.SendError(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}

	err = rc.redisUseCase.Set(req.Context(), key, string(dataJSON), time.Duration(reqBody.Minute)*time.Minute)
	if err != nil {
		pkg.SendError(w, "Failed to store code", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, key)
}

// ==== Router Setup ====
func (rc routerClass) GetClassRouter(r *Router) {
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAllClass))).Methods("GET")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.createClass))).Methods("POST")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.updateClass))).Methods("PATCH")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.deleteClass))).Methods("DELETE")
	r.Router.Handle("/getclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAllClassByEmail))).Methods("GET")
	r.Router.Handle("/class/codeclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.createCodeClass))).Methods("POST")
	r.Router.Handle("/class/joinclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.joinClass))).Methods("POST")
}
