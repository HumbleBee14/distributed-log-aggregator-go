package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/HumbleBee14/distributed-log-aggregator/internal/config"
	"github.com/HumbleBee14/distributed-log-aggregator/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisStorage implements storage operations using Redis
type RedisStorage struct {
	client *redis.Client
	config *config.RedisConfig
	ctx    context.Context
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(cfg *config.RedisConfig) (*RedisStorage, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.ConnectTimeout,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisStorage{
		client: client,
		config: cfg,
		ctx:    ctx,
	}, nil
}

func getServiceKey(serviceName string) string {
	return fmt.Sprintf("service:%s", serviceName)
}

func getLogKey(id string) string {
	return fmt.Sprintf("log:%s", id)
}

// StoreLog stores a log entry in Redis
func (s *RedisStorage) StoreLog(entry models.LogEntry) (string, error) {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	logKey := getLogKey(entry.ID)
	serviceKey := getServiceKey(entry.ServiceName)

	// Create pipeline for atomic operations
	pipe := s.client.Pipeline()

	// Store log entry as hash
	pipe.HSet(s.ctx, logKey, map[string]interface{}{
		"id":           entry.ID,
		"service_name": entry.ServiceName,
		"timestamp":    entry.Timestamp.Format(time.RFC3339),
		"message":      entry.Message,
		"created_at":   entry.CreatedAt.Format(time.RFC3339),
	})

	// Set expiration on the log entry
	pipe.Expire(s.ctx, logKey, s.config.LogExpiryTime)

	// Add to the sorted set with timestamp as score for efficient range queries
	timestampScore := float64(entry.Timestamp.UnixNano())
	pipe.ZAdd(s.ctx, serviceKey, &redis.Z{
		Score:  timestampScore,
		Member: entry.ID,
	})

	// Set expiration on the service key as well
	pipe.Expire(s.ctx, serviceKey, s.config.LogExpiryTime)

	// Execute all commands atomically
	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to store log: %w", err)
	}

	return entry.ID, nil
}

// QueryLogs retrieves logs for a service within a time range
func (s *RedisStorage) QueryLogs(serviceName string, startTime, endTime *time.Time) ([]models.LogResponse, error) {
	serviceKey := getServiceKey(serviceName)

	// Convert time to scores
	var min, max string
	if startTime != nil {
		min = fmt.Sprintf("%d", startTime.UnixNano())
	} else {
		min = "-inf"
	}

	if endTime != nil {
		max = fmt.Sprintf("%d", endTime.UnixNano())
	} else {
		max = "+inf"
	}

	// Get log IDs from the sorted set
	logIDs, err := s.client.ZRangeByScore(s.ctx, serviceKey, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}

	if len(logIDs) == 0 {
		return []models.LogResponse{}, nil
	}

	// Use pipelining to get all log entries efficiently
	pipe := s.client.Pipeline()
	cmds := make(map[string]*redis.StringStringMapCmd)

	for _, id := range logIDs {
		cmds[id] = pipe.HGetAll(s.ctx, getLogKey(id))
	}

	_, err = pipe.Exec(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %w", err)
	}

	// Process results
	logs := make([]models.LogResponse, 0, len(logIDs))

	for _, id := range logIDs {
		result, err := cmds[id].Result()
		if err != nil || len(result) == 0 {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, result["timestamp"])
		if err != nil {
			continue
		}

		logs = append(logs, models.LogResponse{
			Timestamp: timestamp,
			Message:   result["message"],
		})
	}

	// Ensure logs are sorted by timestamp (they should already be, but just in case)
	// We'll implement a proper sort if needed

	return logs, nil
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
