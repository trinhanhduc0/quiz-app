// usecase/auth_usecase.go
package service

import (
	"context"
	"errors"
	"fmt"
	entity "quiz-app/internal/domain/entities"

	"github.com/golang-jwt/jwt/v4"
)

type AuthUseCase struct {
	authService *AuthService
	tokenRepo   RedisUseCase
}

func NewAuthUseCase(authService *AuthService, redisUseCase RedisUseCase) *AuthUseCase {
	return &AuthUseCase{authService: authService, tokenRepo: redisUseCase}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, tokenString string) (entity.AuthClaims, error) {
	// Validate JWT
	token, err := uc.authService.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return entity.AuthClaims{}, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return entity.AuthClaims{}, errors.New("failed to get claims")
	}

	emailID := claims["email_id"].(string)
	redisToken, err := uc.tokenRepo.Get(ctx, fmt.Sprintf("user_token:%s", emailID))
	if err != nil || redisToken != tokenString {
		return entity.AuthClaims{}, errors.New("session expired or user logged in elsewhere")
	}

	return entity.AuthClaims{
		EmailID: claims["email_id"].(string),
		Email:   claims["email"].(string),
		Exp:     int64(claims["exp"].(float64)),
	}, nil
}
