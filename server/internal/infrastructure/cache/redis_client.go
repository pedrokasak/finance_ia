package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache is the port for the caching layer — swap Redis for Memcached without touching callers
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

// RedisCache implements Cache using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache from REDIS_URL env
func NewRedisCache() (*RedisCache, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis: invalid URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: connection failed: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("cache miss: %s", key)
	}
	return val, err
}

func (r *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// NoopCache is a fallback when Redis is unavailable — all operations succeed silently
type NoopCache struct{}

func (n *NoopCache) Get(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("cache miss: noop")
}
func (n *NoopCache) Set(_ context.Context, _, _ string, _ time.Duration) error { return nil }
func (n *NoopCache) Del(_ context.Context, _ ...string) error                  { return nil }
func (n *NoopCache) Exists(_ context.Context, _ string) (bool, error)          { return false, nil }
func (n *NoopCache) Incr(_ context.Context, _ string) (int64, error)           { return 0, nil }
func (n *NoopCache) Expire(_ context.Context, _ string, _ time.Duration) error { return nil }

// CacheKeys — centralized key builders
const (
	KeyDashboard         = "dashboard:%s"     // userID
	KeyAIInsight         = "ai:insight:%s:%s" // userID, period
	KeyAIRateLimit       = "ai:rate:%s:%s"    // userID, weekKey
	KeyCategoriesDefault = "categories:default"
	KeyIdempotency       = "idempotency:%s" // key
)

// TTL constants
const (
	TTLDashboard        = 5 * time.Minute
	TTLAIInsightFree    = 7 * 24 * time.Hour
	TTLAIInsightPremium = 24 * time.Hour
	TTLCategories       = time.Hour
	TTLIdempotency      = 24 * time.Hour
	TTLRateLimit        = 7 * 24 * time.Hour
)
