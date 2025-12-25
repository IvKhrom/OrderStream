package kafkastorage

import "context"

type Message struct {
	Value []byte
}

type Reader interface {
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}

type ReaderFactory interface {
	NewReader(brokers []string, groupID, topic string) Reader
}
