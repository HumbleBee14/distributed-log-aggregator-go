package config

import (
	"os"
	"strconv"
	"time"

	"github.com/HumbleBee14/distributed-log-aggregator/pkg/logger"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Server ServerConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type RedisConfig struct {
	Host           string
	Port           string
	Password       string
	Username       string
	DB             int
	LogExpiryTime  time.Duration
	PoolSize       int
	ConnectTimeout time.Duration
}

// New creates a new configuration with the below precedence:
// 1. System environment variables
// 2. .env file variables
// 3. Default values
func New() *Config {
	// Load .env file if it exists
	// _ = godotenv.Load()
	envFileLoaded := false
	if err := godotenv.Load(); err == nil {
		envFileLoaded = true
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Redis: RedisConfig{
			Host:           getEnv("REDIS_HOST", "localhost"),
			Port:           getEnv("REDIS_PORT", "6379"),
			Password:       getEnv("REDIS_PASSWORD", ""),
			Username:       getEnv("REDIS_USERNAME", "default"),
			DB:             getEnvAsInt("REDIS_DB", 0),
			LogExpiryTime:  getEnvAsDuration("LOG_EXPIRY_TIME", time.Hour),
			PoolSize:       getEnvAsInt("REDIS_POOL_SIZE", 10),
			ConnectTimeout: getEnvAsDuration("REDIS_CONNECT_TIMEOUT", 5*time.Second),
		},
	}

	if !envFileLoaded {
		logger.Warn("No .env file found, using environment variables and defaults")
	} else {
		logger.Info("Configuration loaded from .env")
	}

	logger.Info("Using Redis at %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
