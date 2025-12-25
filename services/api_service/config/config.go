package config

import (
	"os"
	"time"
)

type Config struct {
	PostgresDSN       string
	KafkaBrokers      string
	RedisAddr         string
	HTTPPort          string
	OrdersEventsTopic string
	OrdersAckTopic    string
	AckWaitTimeout    time.Duration
}

type Loader interface {
	Load() (*Config, error)
}

type EnvLoader struct{}

func (EnvLoader) Load() (*Config, error) {
	return &Config{
		PostgresDSN:       getenv("API_POSTGRES_DSN", "postgres://postgres:upvel123@localhost:5433/orderstream_api?sslmode=disable"),
		KafkaBrokers:      getenv("KAFKA_BROKERS", "localhost:29092"),
		RedisAddr:         getenv("REDIS_ADDR", "localhost:6379"),
		HTTPPort:          getenv("API_PORT", "8080"),
		OrdersEventsTopic: getenv("ORDERS_EVENTS_TOPIC", "orders.events"),
		OrdersAckTopic:    getenv("ORDERS_ACK_TOPIC", "orders.ack"),
		AckWaitTimeout:    30 * time.Second,
	}, nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}


