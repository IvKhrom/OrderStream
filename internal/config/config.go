package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN  string
	KafkaBrokers string
	ApiPort      string
	WorkerGroup  string
	RedisAddr    string
}

func Load() (*Config, error) {
	_ = godotenv.Load("config/.env")
	_ = godotenv.Load()

	cfg := &Config{
		PostgresDSN:  getEnv("POSTGRES_DSN", "postgres://postgres:значитбезпароля@localhost:5433/orderstream?sslmode=disable"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		ApiPort:      getEnv("API_PORT", "8080"),
		WorkerGroup:  getEnv("WORKER_GROUP", "order-workers"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
