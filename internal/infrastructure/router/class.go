package routes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/pkg"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type routerClass struct {
	auth         service.AuthHandler
	classUseCase service.ClassUseCase
	redisUseCase service.RedisUseCase
}

func NewRouterClass(s service.ClassUseCase, r service.RedisUseCase, auth service.AuthHandler) routerClass {
	return routerClass{
		auth:         auth,
		classUseCase: s,
		redisUseCase: r,
	}
}

func (rc routerClass) createClass(w http.ResponseWriter, req *http.Request) {
	fmt.Println("CREATE")
	emailID, ok := req.Context().Value("email_id").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}
	email, ok := req.Context().Value("email").(string)
	if !ok {
		pkg.SendError(w, "Invalid email", http.StatusBadRequest)
		return
	}

	var newClass entity.Class
	if err := json.NewDecoder(req.Body).Decode(&newClass); err != nil {
		pkg.SendError(w, "Invalid class field", http.StatusBadRequest)
		return
	}

	// Set metadata fields
	newClass.CreatedAt = time.Now()
	newClass.UpdatedAt = time.Now()
	newClass.AuthorMail = email
	newClass.EmailID = emailID

	// Call use case to create class
	classCreated, err := rc.classUseCase.CreateClass(&newClass)
	if err != nil {
		pkg.SendError(w, "Failed to create class", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, classCreated)
}

func (rc routerClass) getAllClass(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value("email_id").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	// Call use case to get all classes
	allClasses, err := rc.classUseCase.GetAllClass(context.TODO(), emailID)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, allClasses)
}

func (rc routerClass) getAllClassByEmail(w http.ResponseWriter, req *http.Request) {
	email, ok := req.Context().Value("email").(string)
	if !ok {
		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	// Call use case to get all classes
	allClasses, err := rc.classUseCase.GetAllClassByEmail(context.TODO(), email)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, allClasses)
}

func (rc routerClass) updateClass(w http.ResponseWriter, req *http.Request) {
	email := req.Context().Value("email").(string)
	emailID := req.Context().Value("email_id").(string)

	var classUpdate entity.Class
	if err := json.NewDecoder(req.Body).Decode(&classUpdate); err != nil {
		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	classUpdate.EmailID = emailID
	classUpdate.AuthorMail = email
	// Call use case to update class
	updatedClass, err := rc.classUseCase.UpdateClass(context.TODO(), &classUpdate)
	if err != nil {
		fmt.Println(err)
		pkg.SendError(w, "Failed to update class", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, updatedClass)
}

// Delete a class
func (rc routerClass) deleteClass(w http.ResponseWriter, req *http.Request) {
	emailID := req.Context().Value("email_id").(string)

	var classDelete struct {
		ID primitive.ObjectID `json:"_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&classDelete); err != nil {
		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	err := rc.classUseCase.DeleteClass(context.TODO(), emailID, classDelete.ID)
	if err != nil {
		pkg.SendError(w, "Error deleting class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, fmt.Sprintf("Class with ID %v deleted", classDelete.ID))
}

// Join a class
func (rc routerClass) joinClass(w http.ResponseWriter, req *http.Request) {
	email := req.Context().Value("email").(string)

	var class struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&class); err != nil {
		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Get class ID from Redis
	sClassId, err := rc.redisUseCase.Get(req.Context(), class.ID)
	if err != nil {
		pkg.SendError(w, "Error fetching class from Redis", http.StatusInternalServerError)
		return
	}

	// Convert the class ID from string to ObjectID
	oClassId, err := primitive.ObjectIDFromHex(sClassId)
	if err != nil {
		pkg.SendError(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Join the class
	err = rc.classUseCase.JoinClass(context.TODO(), oClassId, email)
	if err != nil {
		fmt.Println(err)
		pkg.SendError(w, "Error joining class", http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully joined class: %s", class.ID),
	}
	pkg.SendResponse(w, http.StatusOK, response)
}

// generateRandomKey generates a random 16-byte string for use as a Redis key
func generateRandomKey() (string, error) {
	bytes := make([]byte, 5) // 16 bytes = 128-bit random key
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (rc routerClass) createCodeClass(w http.ResponseWriter, req *http.Request) {
	var classID struct {
		ID     string `json:"_id"`
		Minute int    `json:"minute"`
	}
	fmt.Println(classID)
	if err := json.NewDecoder(req.Body).Decode(&classID); err != nil {
		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	// Generate a random key
	randomKey, err := generateRandomKey()
	if err != nil {
		pkg.SendError(w, "Failed to generate random key", http.StatusInternalServerError)
		return
	}
	// Store the classID.ID in Redis with the random key and expiration time
	err = rc.redisUseCase.Set(req.Context(), randomKey, classID.ID, time.Minute*time.Duration(classID.Minute))
	if err != nil {
		pkg.SendError(w, "Failed to store data in Redis", http.StatusInternalServerError)
		return
	}
	fmt.Println(rc.redisUseCase.Get(req.Context(), randomKey))
	// Return the random key as a response
	pkg.SendResponse(w, http.StatusOK, randomKey)
}

// GetRouter sets up all routes for class-related operations
func (rc routerClass) GetClassRouter(r *Router) {
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAllClass))).Methods("GET")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.createClass))).Methods("POST")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.updateClass))).Methods("PATCH")
	r.Router.Handle("/class", rc.auth.AuthMiddleware(http.HandlerFunc(rc.deleteClass))).Methods("DELETE")
	r.Router.Handle("/getclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.getAllClassByEmail))).Methods("GET")

	r.Router.Handle("/class/codeclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.createCodeClass))).Methods("POST")
	r.Router.Handle("/class/joinclass", rc.auth.AuthMiddleware(http.HandlerFunc(rc.joinClass))).Methods("POST")
}
