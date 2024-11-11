package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	"quiz-app/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps a Redis client with additional functionality
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client and checks the connection
func NewRedisClient(redisURL string) (repository.RedisRepository, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

// Basic Redis operations

// Set adds a new key-value pair in Redis
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %v", key, err)
	}
	return nil
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %v", key, err)
	}
	return val, nil
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", key, err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %v", key, err)
	}
	return count > 0, nil
}

// Close closes the Redis client connection
func (r *RedisClient) Close(ctx context.Context) {
	r.client.Close()
}

// List operations

func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.LPush(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("failed to push to list %s: %v", key, err)
	}
	return nil
}

func (r *RedisClient) RPop(ctx context.Context, key string) (string, error) {
	val, err := r.client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to pop from list %s: %v", key, err)
	}
	return val, nil
}

func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	values, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list %s: %v", key, err)
	}
	return values, nil
}

// Set operations

func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to set %s: %v", key, err)
	}
	return nil
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get set %s: %v", key, err)
	}
	return members, nil
}

// Sorted Set operations

func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	err := r.client.ZAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to sorted set %s: %v", key, err)
	}
	return nil
}

func (r *RedisClient) ZRangeByScore(ctx context.Context, key, min, max string) ([]string, error) {
	members, err := r.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get sorted set %s: %v", key, err)
	}
	return members, nil
}

// Hash operations

func (r *RedisClient) HSet(ctx context.Context, key string, fields map[string]interface{}) error {
	err := r.client.HMSet(ctx, key, fields).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash %s: %v", key, err)
	}
	return nil
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash %s: %v", key, err)
	}
	return values, nil
}

// Geospatial operations

func (r *RedisClient) GeoAdd(ctx context.Context, key string, locations ...*redis.GeoLocation) error {
	err := r.client.GeoAdd(ctx, key, locations...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to geo %s: %v", key, err)
	}
	return nil
}

func (r *RedisClient) GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error) {
	positions, err := r.client.GeoPos(ctx, key, members...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get geo positions for %s: %v", key, err)
	}
	return positions, nil
}

// GetRedis creates and returns a Redis client instance
func GetRedis() (repository.RedisRepository, error) {
	// Upstash Redis URL - thay bằng URL của bạn
	redisURL := os.Getenv("REDIS_URI") // This is where you fetch the secret

	//redisURL := "rediss://default:YOUR_UPSTASH_PASSWORD@YOUR_UPSTASH_ENDPOINT:6379"
	return NewRedisClient(redisURL)
}
