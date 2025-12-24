package kafkastorage

import "context"

// Message — минимальное сообщение, которое нужно consumer-слою (Value).
// Без kafka.Message в consumer, чтобы слой не зависел от kafka-go.
type Message struct {
	Value []byte
}

// Reader читает сообщения из Kafka.
type Reader interface {
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}

// ReaderFactory создаёт Reader для заданных параметров.
type ReaderFactory interface {
	NewReader(brokers []string, groupID, topic string) Reader
}
