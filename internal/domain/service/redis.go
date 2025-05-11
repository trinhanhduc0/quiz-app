package service

import (
	"context"
	"quiz-app/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUseCase struct {
	RedisRepo repository.RedisRepository
}

// NewRedisUseCase initializes a new RedisUseCase with the given Redis repository
func NewRedisUseCase(redisRepo repository.RedisRepository) *RedisUseCase {
	return &RedisUseCase{
		RedisRepo: redisRepo,
	}
}

// Basic Redis operations

// Set stores a key-value pair in Redis with an expiration time
func (uc *RedisUseCase) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return uc.RedisRepo.Set(ctx, key, value, expiration)
}

// Get retrieves the value associated with a key from Redis
func (uc *RedisUseCase) Get(ctx context.Context, key string) (string, error) {
	return uc.RedisRepo.Get(ctx, key)
}

// Delete removes a key from Redis
func (uc *RedisUseCase) Delete(ctx context.Context, key string) error {
	return uc.RedisRepo.Delete(ctx, key)
}

// Exists checks if a key exists in Redis
func (uc *RedisUseCase) Exists(ctx context.Context, key string) (bool, error) {
	return uc.RedisRepo.Exists(ctx, key)
}

// List operations

// LPush pushes values to the left of a list in Redis
func (uc *RedisUseCase) LPush(ctx context.Context, key string, expiration time.Duration, values ...interface{}) error {
	return uc.RedisRepo.LPush(ctx, key, expiration, values...)
}

// RPop removes and retrieves the last element from a list in Redis
func (uc *RedisUseCase) RPop(ctx context.Context, key string) (string, error) {
	return uc.RedisRepo.RPop(ctx, key)
}

// LRange retrieves a range of elements from a list in Redis
func (uc *RedisUseCase) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return uc.RedisRepo.LRange(ctx, key, start, stop)
}

// Set operations

// SAdd adds members to a set in Redis
func (uc *RedisUseCase) SAdd(ctx context.Context, key string, expiration time.Duration, members ...interface{}) error {
	return uc.RedisRepo.SAdd(ctx, key, expiration, members...)
}

// SMembers retrieves all members of a set in Redis
func (uc *RedisUseCase) SMembers(ctx context.Context, key string) ([]string, error) {
	return uc.RedisRepo.SMembers(ctx, key)
}

// Sorted Set operations

// ZAdd adds elements with scores to a sorted set in Redis
func (uc *RedisUseCase) ZAdd(ctx context.Context, key string, expiration time.Duration, members ...redis.Z) error {
	return uc.RedisRepo.ZAdd(ctx, key, expiration, members...)
}

// ZRangeByScore retrieves elements from a sorted set by score range in Redis
func (uc *RedisUseCase) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	return uc.RedisRepo.ZRangeByScore(ctx, key, min, max)
}

// Hash operations

// HSet sets multiple fields in a hash in Redis
func (uc *RedisUseCase) HSet(ctx context.Context, key string, expiration time.Duration, values map[string]interface{}) error {
	return uc.RedisRepo.HSet(ctx, key, expiration, values)
}

// HGetAll retrieves all fields and values from a hash in Redis
func (uc *RedisUseCase) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return uc.RedisRepo.HGetAll(ctx, key)
}

// Geospatial operations

// GeoAdd adds a geospatial location to a key in Redis
func (uc *RedisUseCase) GeoAdd(ctx context.Context, key string, expiration time.Duration, location *redis.GeoLocation) error {
	return uc.RedisRepo.GeoAdd(ctx, key, expiration, location)
}

// GeoPos retrieves the geospatial position of specified members from a key in Redis
func (uc *RedisUseCase) GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error) {
	return uc.RedisRepo.GeoPos(ctx, key, members...)
}

// Close closes the Redis connection
func (uc *RedisUseCase) Close(ctx context.Context) {
	uc.RedisRepo.Close(ctx)
}

func (uc *RedisUseCase) Lock(ctx context.Context, key string) error {
	return uc.RedisRepo.Lock(ctx, key)
}

func (uc *RedisUseCase) Unlock(ctx context.Context, key string) {
	uc.RedisRepo.Unlock(ctx, key)
}
