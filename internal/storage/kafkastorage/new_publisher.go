package kafkastorage

import "github.com/segmentio/kafka-go"

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}


