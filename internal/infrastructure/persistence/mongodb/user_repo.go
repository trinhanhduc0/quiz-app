package persistence

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserMongoRepository implements the repository.UserRepository interface
type UserMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

// NewUserMongoRepository tạo một instance mới của UserMongoRepository
func NewUserMongoRepository() repository.UserRepository {
	collRepo := NewCollRepository("userdb", "users")
	return &UserMongoRepository{
		CollRepo: collRepo,
	}
}

func (ur UserMongoRepository) CreateUser(ctx context.Context, user *entity.User) (primitive.ObjectID, error) {
	var errHash error
	user.Password, errHash = hashPassword(user.Password)
	if errHash != nil {
		return primitive.NilObjectID, errHash
	}
	newUserID, err := ur.CollRepo.Create(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	idUser := newUserID.(primitive.ObjectID)
	return idUser, nil
}

// HashPassword generates a bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (ur UserMongoRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	filter := bson.M{"email_id": user.EmailID}
	ur.CollRepo.Update(ctx, filter, user)
	return &entity.User{}, nil
}

func (ur UserMongoRepository) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
	return &entity.User{}, nil
}

func (ur UserMongoRepository) GetUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	filter := bson.M{"email_id": user.EmailID}
	userInfo, err := ur.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if userInfo is a map
	userData, ok := userInfo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for userInfo")
	}

	// Map the fields to the user entity
	userEntity := &entity.User{
		ID:      userData["_id"].(primitive.ObjectID),
		EmailID: userData["email_id"].(string),
	}

	return userEntity, nil
}

func (ur UserMongoRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return fmt.Errorf("failed to convert inserted ID to ObjectID")
}
