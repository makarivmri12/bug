package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// RedisCache handles caching and task queue operations
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client, logger *zap.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
	}
}

// PushTask adds a task to the queue
func (rc *RedisCache) PushTask(ctx context.Context, queue string, task interface{}) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		rc.logger.Error("failed to marshal task", zap.Error(err))
		return err
	}

	err = rc.client.RPush(ctx, queue, taskJSON).Err()
	if err != nil {
		rc.logger.Error("failed to push task", zap.Error(err))
		return err
	}

	return nil
}

// PopTask retrieves and removes a task from the queue
func (rc *RedisCache) PopTask(ctx context.Context, queue string, timeout time.Duration) (string, error) {
	result, err := rc.client.BLPop(ctx, timeout, queue).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Timeout
		}
		rc.logger.Error("failed to pop task", zap.Error(err))
		return "", err
	}

	if len(result) < 2 {
		return "", fmt.Errorf("invalid pop result")
	}

	return result[1], nil
}

// SetScanMetadata stores scan metadata with TTL
func (rc *RedisCache) SetScanMetadata(ctx context.Context, scanID string, metadata map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("scan:%s", scanID)
	data, err := json.Marshal(metadata)
	if err != nil {
		rc.logger.Error("failed to marshal metadata", zap.Error(err))
		return err
	}

	err = rc.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		rc.logger.Error("failed to set scan metadata", zap.Error(err))
		return err
	}

	return nil
}

// GetScanMetadata retrieves scan metadata
func (rc *RedisCache) GetScanMetadata(ctx context.Context, scanID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("scan:%s", scanID)

	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		rc.logger.Error("failed to get scan metadata", zap.Error(err))
		return nil, err
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(val), &metadata)
	if err != nil {
		rc.logger.Error("failed to unmarshal metadata", zap.Error(err))
		return nil, err
	}

	return metadata, nil
}

// DeleteScanMetadata removes scan metadata
func (rc *RedisCache) DeleteScanMetadata(ctx context.Context, scanID string) error {
	key := fmt.Sprintf("scan:%s", scanID)

	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		rc.logger.Error("failed to delete scan metadata", zap.Error(err))
		return err
	}

	return nil
}

// SetRateLimit sets a rate limit key
func (rc *RedisCache) SetRateLimit(ctx context.Context, key string, limit int, ttl time.Duration) error {
	return rc.client.Set(ctx, fmt.Sprintf("rate:%s", key), limit, ttl).Err()
}

// CheckRateLimit checks if rate limit is exceeded
func (rc *RedisCache) CheckRateLimit(ctx context.Context, key string) (int, error) {
	val, err := rc.client.Get(ctx, fmt.Sprintf("rate:%s", key)).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}

// IncrementRateLimit increments the rate limit counter
func (rc *RedisCache) IncrementRateLimit(ctx context.Context, key string) (int, error) {
	val, err := rc.client.Incr(ctx, fmt.Sprintf("rate:%s", key)).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
