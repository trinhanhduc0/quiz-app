package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
)

var (
	firebaseConfigFile = "firebase-config.json"
	authClient         *auth.Client
	once               sync.Once
	jwtSecret          = []byte("serect_key")
)

// GetAuth returns a singleton Firebase Auth client
func GetAuth() (*auth.Client, error) {
	var err error
	once.Do(func() {
		// Initialize Firebase
		ctx := context.Background()
		opt := option.WithCredentialsFile(firebaseConfigFile)

		app, initErr := firebase.NewApp(ctx, nil, opt)
		if initErr != nil {
			log.Fatalf("Error initializing Firebase app: %v", initErr)
			err = initErr
			return
		}

		authClient, err = app.Auth(ctx)
		if err != nil {
			log.Fatalf("Error initializing Firebase Auth client: %v", err)
		}
	})

	return authClient, err
}

func CreateJWT(user_id primitive.ObjectID, acc_id, email string) (string, error) {
	claims := jwt.MapClaims{
		"_id":      user_id,
		"email_id": acc_id,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}
	return jwtToken, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func CheckTokenExpiry(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to get claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("failed to get exp claim")
	}

	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf("token has expired")
	}

	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		err = CheckTokenExpiry(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Token expired: %v", err), http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["_id"])
		ctx := context.WithValue(r.Context(), "_id", claims["_id"])
		ctx = context.WithValue(ctx, "email_id", claims["email_id"])
		ctx = context.WithValue(ctx, "email", claims["email"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
