package service

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassUseCase struct {
	repoClass repository.ClassRepository
	repoTest  repository.TestRepository
}

func NewClassUseCase(repoClass repository.ClassRepository, repoTest repository.TestRepository) *ClassUseCase {
	return &ClassUseCase{repoClass: repoClass, repoTest: repoTest}
}

func (cuc *ClassUseCase) CreateClass(class *entity.Class) (any, error) {
	// Step 1: Create class in the database
	newID, err := cuc.repoClass.CreateClass(context.TODO(), class)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Error creating class: %v\n", err)
		return nil, fmt.Errorf("failed to create class: %w", err)
	}

	// Set the new class ID
	class.ID = newID

	// Step 2: Update allowed users for the associated test
	err = cuc.repoTest.UpdateAllowUser(context.TODO(), class.TestID, class.StudentAccept)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Error updating allowed users for test %v: %v\n", class.TestID, err)
		return nil, fmt.Errorf("failed to update allowed users for test %v: %w", class.TestID, err)
	}

	// If everything succeeds, return the class
	return class, nil
}

func (cuc *ClassUseCase) GetAllClassByEmail(ctx context.Context, email string) ([]any, error) {
	return cuc.repoClass.GetAllClassByEmail(context.TODO(), email)
}

func (cuc *ClassUseCase) UpdateClass(ctx context.Context, class *entity.Class) (any, error) {
	// Step 2: Update allowed users for the associated test
	err := cuc.repoTest.UpdateAllowUser(context.TODO(), class.TestID, class.StudentAccept)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Error updating allowed users for test %v: %v\n", class.TestID, err)
		return nil, fmt.Errorf("failed to update allowed users for test %v: %w", class.TestID, err)
	}

	return cuc.repoClass.UpdateClass(context.TODO(), class)
}

func (cuc *ClassUseCase) DeleteClass(ctx context.Context, emailID string, id primitive.ObjectID) error {
	return cuc.repoClass.DeleteClass(context.TODO(), emailID, id)
}

func (cuc *ClassUseCase) GetAllClass(ctx context.Context, id string) ([]any, error) {
	return cuc.repoClass.GetClassByAuthorEmail(context.TODO(), id)
}

func (cuc *ClassUseCase) JoinClass(ctx context.Context, classID primitive.ObjectID, email string) error {
	return cuc.repoClass.JoinClass(ctx, classID, email)
}
