package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

//go:generate mockery --name ProducerClient --output ../mocks --outpkg mocks --case underscore
//go:generate mockery --name ConsumerClient --output ../mocks --outpkg mocks --case underscore

// Публикация в Kafka
type ProducerClient interface {
	Publish(ctx context.Context, value []byte) error
	Close() error
}

// Consumer в Kafka
type ConsumerClient interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}
