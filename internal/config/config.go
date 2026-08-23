package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, LogLevel string
	SessionTTL, WorkerInterval      time.Duration
	Clock                           func() time.Time
}

func Load() Config {
	return Config{HTTPAddr: env("HTTP_ADDR", ":8080"), DatabaseURL: env("DATABASE_URL", "file:windsea.db?_pragma=foreign_keys(1)"), LogLevel: env("LOG_LEVEL", "INFO"), SessionTTL: duration("SESSION_TTL", 24*time.Hour), WorkerInterval: duration("WORKER_INTERVAL", 2*time.Second), Clock: time.Now}
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func duration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}
func Bool(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return value
}
