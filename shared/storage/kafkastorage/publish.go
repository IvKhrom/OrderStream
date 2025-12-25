package kafkastorage

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func (p *Publisher) Publish(ctx context.Context, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{
		Value: value,
		Time:  time.Now(),
	})
}


