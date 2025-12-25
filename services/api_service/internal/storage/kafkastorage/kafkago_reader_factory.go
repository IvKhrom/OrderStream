package kafkastorage

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaGoReader struct {
	r *kafka.Reader
}

func (kr *kafkaGoReader) ReadMessage(ctx context.Context) (Message, error) {
	msg, err := kr.r.ReadMessage(ctx)
	if err != nil {
		return Message{}, err
	}
	return Message{Value: msg.Value}, nil
}

func (kr *kafkaGoReader) Close() error {
	return kr.r.Close()
}

type KafkaGoReaderFactory struct {
	HeartbeatInterval time.Duration
	SessionTimeout    time.Duration
}

func NewKafkaGoReaderFactory() *KafkaGoReaderFactory {
	return &KafkaGoReaderFactory{
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	}
}

func (f *KafkaGoReaderFactory) NewReader(brokers []string, groupID, topic string) Reader {
	return &kafkaGoReader{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers:           brokers,
			GroupID:           groupID,
			Topic:             topic,
			HeartbeatInterval: f.HeartbeatInterval,
			SessionTimeout:    f.SessionTimeout,
		}),
	}
}
