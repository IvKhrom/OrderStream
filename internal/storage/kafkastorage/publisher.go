package kafkastorage

import (
	"context"
	"encoding/json"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/segmentio/kafka-go"
)

// writer — минимальный интерфейс для kafka.Writer, чтобы kafkastorage можно было тестировать без реальной Kafka.
// Интерфейс используется только внутри storage слоя.
type writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Publisher struct {
	w writer
}

func (p *Publisher) PublishOrderEvent(ctx context.Context, event *models.OrderEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.Publish(ctx, b)
}

func (p *Publisher) PublishOrderAck(ctx context.Context, ack *models.OrderAck) error {
	b, err := json.Marshal(ack)
	if err != nil {
		return err
	}
	return p.Publish(ctx, b)
}
