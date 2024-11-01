package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	persistence "quiz-app/internal/infrastructure/persistence/mongodb"

	"firebase.google.com/go/v4/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesAuth struct {
	r *Router
}

func GetRoutesAuth(r *Router) *RoutesAuth {
	return &RoutesAuth{
		r: r,
	}
}

// Login handles Google login using Firebase Auth
func (ra *RoutesAuth) handlerLogin(w http.ResponseWriter, r *http.Request, authClient *auth.Client, ctx context.Context) {
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
	verifyTokenResp, err := authClient.VerifyIDToken(ctx, req.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error verifying token: %v", err), http.StatusInternalServerError)
		return
	}

	emailID := verifyTokenResp.UID
	emailUser, ok := verifyTokenResp.Claims["email"].(string)
	if !ok {
		http.Error(w, "Error extracting email from token", http.StatusInternalServerError)
		return
	}

	repoUser := persistence.NewUserMongoRepository()
	userRow, err := repoUser.GetUser(context.TODO(), &entity.User{EmailID: emailID})
	// // Define user data
	// user, err := entity.NewUser(
	// 	emailID, "", "", emailUser, "default_password",
	// )

	var userID primitive.ObjectID
	if userRow == nil {
		// // Create new user if they don't exist
		// userID, err = repoUser.Create(user)
		// if err != nil {
		// 	http.Error(w, "Failed to create user", http.StatusInternalServerError)
		// 	return
		// }
	} else {
		// Use existing user ID
		userID = userRow.ID
	}

	// Generate JWT token
	token, err := service.CreateJWT(userID, emailID, emailUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create JWT: %v", err), http.StatusInternalServerError)
		return
	}

	// Send JWT token in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (r *RoutesAuth) login(w http.ResponseWriter, req *http.Request) {
	authClient, err := service.GetAuth()
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v", err)
		return
	}
	// Add custom handler for Google login to the router
	r.handlerLogin(w, req, authClient, context.TODO())
}

func (ra *RoutesAuth) SetRouteLogin() {
	ra.r.Handle("/api/google/login", http.HandlerFunc(ra.login)).Methods("POST")
}
