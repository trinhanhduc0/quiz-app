// internal/adapter/routes/auth_routes.go

package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/pkg"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesAuth struct {
	router       *Router
	authService  service.AuthService
	userUsecase  service.UserUseCase
	redisUseCase service.RedisUseCase
	classUseCase service.ClassUseCase
}

func NewRoutesAuth(r *Router, authService service.AuthService, userUsecase service.UserUseCase, redisUseCase service.RedisUseCase, classUseCase service.ClassUseCase) *RoutesAuth {
	return &RoutesAuth{
		router:       r,
		authService:  authService,
		userUsecase:  userUsecase,
		redisUseCase: redisUseCase,
		classUseCase: classUseCase,
	}
}

// Login handles Google login using Firebase Auth
func (ra *RoutesAuth) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Missing token parameter", http.StatusBadRequest)
		return
	}

	// Verify Firebase ID token
	ctx := context.TODO()
	verifyTokenResp, err := ra.authService.VerifyIDToken(ctx, req.Token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error verifying token: %v", err), http.StatusInternalServerError)
		return
	}

	emailID := verifyTokenResp.UID
	emailUser, ok := verifyTokenResp.Claims["email"].(string)
	nameUser := strings.Split(verifyTokenResp.Claims["name"].(string), " ")
	if !ok {
		http.Error(w, "Error extracting email from token", http.StatusInternalServerError)
		return
	}

	userRow, err := ra.userUsecase.GetUser(ctx, &entity.User{EmailID: emailID})
	var userID primitive.ObjectID
	if userRow == nil {
		user, err := entity.NewUser(
			emailID, nameUser[len(nameUser)-1], nameUser[0], emailUser, "default_password",
		)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Create new user if they don't exist
		userID, err = ra.userUsecase.CreateUser(ctx, user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// CLASS TEST-----------------------------------------------------------------------------
		// Lấy dữ liệu từ Redis
		dataJSON, err := ra.redisUseCase.Get(r.Context(), "93d433c9be")
		if err != nil {
			pkg.SendError(w, "Error fetching class from Redis", http.StatusInternalServerError)
			return
		}

		// Định nghĩa cấu trúc để giải mã dữ liệu JSON từ Redis
		var data struct {
			ClassID     string   `json:"class_id"`
			EmailAuthor string   `json:"email"`
			TestID      []string `json:"test_id"`
		}

		// Giải mã dữ liệu JSON
		err = json.Unmarshal([]byte(dataJSON), &data)
		if err != nil {
			pkg.SendError(w, "Error decoding class data", http.StatusInternalServerError)
			return
		}

		// Chuyển đổi class ID từ chuỗi sang ObjectID
		oClassId, err := primitive.ObjectIDFromHex(data.ClassID)
		if err != nil {
			pkg.SendError(w, "Invalid class ID", http.StatusBadRequest)
			return
		}
		var oTestId []primitive.ObjectID
		for _, v := range data.TestID {
			// Chuyển đổi từng chuỗi testID sang ObjectID
			testID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				pkg.SendError(w, "Invalid test ID", http.StatusBadRequest)
				return
			}
			// Thêm ObjectID đã chuyển đổi vào slice oTestId
			oTestId = append(oTestId, testID)
		}

		// Thêm user vào lớp học
		err = ra.classUseCase.JoinClass(context.TODO(), oClassId, oTestId, data.EmailAuthor, emailUser)
		if err != nil {
			fmt.Println(err)
			pkg.SendError(w, "Error joining class", http.StatusInternalServerError)
			return
		}

	} else {
		// Use existing user ID
		userID = userRow.ID
	}

	// Generate JWT token
	token, err := ra.authService.CreateJWT(entity.AuthClaims{
		UserID: userID, EmailID: emailID, Email: emailUser, Exp: 24,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create JWT: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the JWT token in Redis with the userID as the key
	err = ra.redisUseCase.Set(ctx, fmt.Sprintf("user_token:%s", emailID), token, time.Hour*24)
	if err != nil {
		http.Error(w, "Failed to store token in Redis", http.StatusInternalServerError)
		return
	}

	// Send JWT token in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// SetLoginRoute registers the login route with the router
func (ra *RoutesAuth) SetLoginRoute() {
	ra.router.Handle("/api/google/login", http.HandlerFunc(ra.loginHandler)).Methods("POST")
}
