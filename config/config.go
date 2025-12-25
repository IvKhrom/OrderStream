package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN       string
	KafkaBrokers      string
	ApiPort           string
	WorkerGroup       string
	OrdersEventsTopic string
	OrdersAckTopic    string
	AckWaitTimeout    time.Duration
}

func LoadConfig(_ string) (*Config, error) {
	_ = godotenv.Overload("config/.env")
	_ = godotenv.Overload()

	cfg := &Config{
		PostgresDSN:       getEnv("POSTGRES_DSN", "postgres://postgres:upvel123@localhost:5433/orderstream?sslmode=disable"),
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
