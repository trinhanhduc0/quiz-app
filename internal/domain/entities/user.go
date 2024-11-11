package entity

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	EmailID   string             `bson:"email_id" json:"email_id"`
}


type AuthClaims struct {
	UserID   primitive.ObjectID `json:"_id"`
	EmailID  string             `json:"email_id"`
	Email    string             `json:"email"`
	Exp      int64              `json:"exp"`
}

// NewUser constructs a new User entity and performs validation.
func NewUser(email_id, firstName, lastName, email, password string) (*User, error) {
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return nil, errors.New("first and last name cannot be empty")
	}
	if !isValidEmail(email) {
		return nil, errors.New("invalid email address")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	return &User{
		ID:        primitive.NewObjectID(),
		EmailID:   email_id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     strings.TrimSpace(email),
		Password:  hashPassword(password), // Hash the password
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(), // Set initial updated_at time
	}, nil
}

// UpdatePassword allows updating the user's password with hashing.
func (u *User) UpdatePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	u.Password = hashPassword(newPassword)
	u.UpdatedAt = time.Now() // Update the updated_at field
	return nil
}

// Auxiliary functions for validation and hashing

// Placeholder for hashing function
func hashPassword(password string) string {
	// Implement the actual password hashing here
	return password // Replace with a proper hash function
}

// Simple email validation
func isValidEmail(email string) bool {
	// Implement email validation logic (can use regex or external libraries)
	return strings.Contains(email, "@")
}
