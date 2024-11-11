// package service

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	firebase "firebase.google.com/go/v4"
// 	"firebase.google.com/go/v4/auth" // Add the Redis package
// 	"github.com/golang-jwt/jwt/v4"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"google.golang.org/api/option"
// )

// var (
// 	authClient  *auth.Client
// 	redisClient *service.RedisUseCase
// 	once        sync.Once
// 	jwtSecret   = []byte("secret_key")
// )

// // Initialize Redis client
// func initRedisClient() {
// 	once.Do(func() {
// 		redisRepo, err := redisdb.GetRedis()
// 		redisClient = service.NewRedisUseCase(redisRepo)
// 	})
// }

// type FirebaseConfig struct {
// 	Type                string `json:"type"`
// 	ProjectID           string `json:"project_id"`
// 	PrivateKeyID        string `json:"private_key_id"`
// 	PrivateKey          string `json:"private_key"`
// 	ClientEmail         string `json:"client_email"`
// 	ClientID            string `json:"client_id"`
// 	AuthURI             string `json:"auth_uri"`
// 	TokenURI            string `json:"token_uri"`
// 	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
// 	ClientCertURL       string `json:"client_x509_cert_url"`
// }

// func GetAuth() (*auth.Client, error) {
// 	var err error
// 	once.Do(func() {
// 		ctx := context.Background()
// 		opt := option.WithCredentialsFile("firebase-config.json")

// 		app, initErr := firebase.NewApp(ctx, nil, opt)
// 		if initErr != nil {
// 			log.Fatalf("Error initializing Firebase app: %v", initErr)
// 			err = initErr
// 			return
// 		}

// 		authClient, err = app.Auth(ctx)
// 		if err != nil {
// 			log.Fatalf("Error initializing Firebase Auth client: %v", err)
// 		}
// 	})

// 	return authClient, err
// }

// func CreateJWT(user_id primitive.ObjectID, acc_id, email string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"_id":      user_id,
// 		"email_id": acc_id,
// 		"email":    email,
// 		"exp":      time.Now().Add(time.Hour * 24).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	jwtToken, err := token.SignedString(jwtSecret)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to sign JWT token: %v", err)
// 	}
// 	return jwtToken, nil
// }

// func ValidateJWT(tokenString string) (*jwt.Token, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return jwtSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	return token, nil
// }

// func CheckTokenExpiry(token *jwt.Token) error {
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return fmt.Errorf("failed to get claims")
// 	}

// 	exp, ok := claims["exp"].(float64)
// 	if !ok {
// 		return fmt.Errorf("failed to get exp claim")
// 	}

// 	if time.Now().Unix() > int64(exp) {
// 		return fmt.Errorf("token has expired")
// 	}

// 	return nil
// }

// // Middleware for authentication and Redis token validation
// func AuthMiddleware(next http.Handler) http.Handler {
// 	initRedisClient() // Ensure Redis client is initialized
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			http.Error(w, "Missing token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Remove "Bearer " prefix if present
// 		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
// 			tokenString = tokenString[7:]
// 		}

// 		token, err := ValidateJWT(tokenString)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
// 			return
// 		}

// 		err = CheckTokenExpiry(token)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Token expired: %v", err), http.StatusUnauthorized)
// 			return
// 		}

// 		claims, _ := token.Claims.(jwt.MapClaims)
// 		userID := claims["_id"].(string)

// 		// Check Redis to see if the token matches the stored token
// 		redisToken, err := redisClient.Get(r.Context(), fmt.Sprintf("user_token:%s", userID)).Result()
// 		if err != nil || redisToken != tokenString {
// 			http.Error(w, "Session expired or user logged in elsewhere", http.StatusUnauthorized)
// 			return
// 		}

// 		// Store user information in the context for further use
// 		ctx := context.WithValue(r.Context(), "_id", userID)
// 		ctx = context.WithValue(ctx, "email_id", claims["email_id"])
// 		ctx = context.WithValue(ctx, "email", claims["email"])

//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
package service

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/option"
)

type AuthService struct {
	firebaseClient *auth.Client
	jwtSecret      []byte
}

func NewAuthService(firebaseConfigFile string, jwtSecret []byte) (*AuthService, error) {
	// jsonString := os.Getenv("FIREBASE_CONFIG")
	// opt := option.WithCredentialsJSON([]byte(jsonString))

	//opt := option.WithCredentialsFile(firebaseConfigFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	firebaseClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &AuthService{firebaseClient: firebaseClient, jwtSecret: jwtSecret}, nil
}

func (s *AuthService) CreateJWT(claims entity.AuthClaims) (string, error) {
	tokenClaims := jwt.MapClaims{
		"_id":      claims.UserID,
		"email_id": claims.EmailID,
		"email":    claims.Email,
		"exp":      time.Now().Add(time.Hour * time.Duration(claims.Exp)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
}

func (s *AuthService) VerifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.firebaseClient.VerifyIDToken(ctx, token)
}
