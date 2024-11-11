package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RedisObject defines a Redis-friendly object structure.
type RedisObject struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`                // Object identifier
	Name      string             `bson:"name" json:"name"`             // Object name
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Creation timestamp
	ExpiredAt time.Time          `bson:"expired_at" json:"expired_at"` // Expiration timestamp
}

// CreateRedisObject initializes a RedisObject with an expiration duration.
func CreateRedisObject(name string, duration time.Duration) (RedisObject, error) {
	// Validate that the name is not empty.
	if name == "" {
		return RedisObject{}, errors.New("name cannot be empty")
	}

	// Initialize a RedisObject with ID, CreatedAt, and ExpiredAt fields.
	object := RedisObject{
		ID:        primitive.NewObjectID(),
		Name:      name,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return object, nil
}

// HasExpired checks if the RedisObject has expired based on the current time.
func (r RedisObject) HasExpired() bool {
	return time.Now().After(r.ExpiredAt)
}

// TimeToExpire returns the time remaining until expiration.
func (r RedisObject) TimeToExpire() time.Duration {
	return time.Until(r.ExpiredAt)
}

// ToMap converts a RedisObject to a map for Redis storage as a hash.
func (r RedisObject) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         r.ID.Hex(),
		"name":       r.Name,
		"created_at": r.CreatedAt.Format(time.RFC3339), // Store dates in ISO 8601 format.
		"expired_at": r.ExpiredAt.Format(time.RFC3339),
	}
}

// FromMap converts a Redis hash map back into a RedisObject.
func FromMap(data map[string]string) (RedisObject, error) {
	// Parse the object ID from hex string.
	id, err := primitive.ObjectIDFromHex(data["id"])
	if err != nil {
		return RedisObject{}, err
	}

	// Parse the creation and expiration times.
	createdAt, err := time.Parse(time.RFC3339, data["created_at"])
	if err != nil {
		return RedisObject{}, err
	}

	expiredAt, err := time.Parse(time.RFC3339, data["expired_at"])
	if err != nil {
		return RedisObject{}, err
	}

	return RedisObject{
		ID:        id,
		Name:      data["name"],
		CreatedAt: createdAt,
		ExpiredAt: expiredAt,
	}, nil
}
