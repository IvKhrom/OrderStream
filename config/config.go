package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN      string
	KafkaBrokers     string
	ApiPort          string
	WorkerGroup      string
	OrdersEventsTopic string
	OrdersAckTopic    string
	AckWaitTimeout    time.Duration
}

func LoadConfig(_ string) (*Config, error) {
	// Важно: godotenv.Load НЕ перезаписывает уже установленные переменные окружения.
	// Если в системе/терминале когда‑то был задан KAFKA_BROKERS=kafka:9092,
	// то Load оставит его как есть, и локальный запуск будет пытаться резолвить "kafka".
	//
	// Поэтому используем Overload: локальный config/.env должен иметь приоритет при go run.
	_ = godotenv.Overload("config/.env")
	_ = godotenv.Overload()

	cfg := &Config{
		// Для локального запуска с docker-compose (Postgres проброшен на 5433).
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:upvel123@localhost:5433/orderstream?sslmode=disable"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "localhost:29092"),
		ApiPort:           getEnv("API_PORT", "8080"),
		WorkerGroup:       getEnv("WORKER_GROUP", "order-workers"),
		OrdersEventsTopic: getEnv("ORDERS_EVENTS_TOPIC", "orders.events"),
		OrdersAckTopic:    getEnv("ORDERS_ACK_TOPIC", "orders.ack"),
		AckWaitTimeout:    30 * time.Second,
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


