package kafkastorage

import "github.com/segmentio/kafka-go"

type Publisher struct {
	w *kafka.Writer
}


