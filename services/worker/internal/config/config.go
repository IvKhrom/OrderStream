package config

import (
	"os"
)

type Config struct {
	PostgresDSN       string
	KafkaBrokers      string
	RedisAddr         string
	WorkerGroup       string
	OrdersEventsTopic string
	OrdersAckTopic    string
}

type Loader interface {
	Load() (*Config, error)
}

type EnvLoader struct{}

func (EnvLoader) Load() (*Config, error) {
	return &Config{
		PostgresDSN:       getenv("WORKER_POSTGRES_DSN", "postgres://postgres:upvel123@localhost:5434/orderstream_worker?sslmode=disable"),
		KafkaBrokers:      getenv("KAFKA_BROKERS", "localhost:29092"),
		RedisAddr:         getenv("REDIS_ADDR", "localhost:6379"),
		WorkerGroup:       getenv("WORKER_GROUP", "order-workers"),
		OrdersEventsTopic: getenv("ORDERS_EVENTS_TOPIC", "orders.events"),
		OrdersAckTopic:    getenv("ORDERS_ACK_TOPIC", "orders.ack"),
	}, nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}


