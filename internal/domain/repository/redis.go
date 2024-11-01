package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close(ctx context.Context)

	// List operations
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)

	// Sorted Set operations
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error)

	// Hash operations
	HSet(ctx context.Context, key string, values map[string]interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// Geospatial operations
	GeoAdd(ctx context.Context, key string, location ...*redis.GeoLocation) error
	GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error)
}
