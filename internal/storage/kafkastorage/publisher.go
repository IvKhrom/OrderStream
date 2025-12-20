package kafkastorage

import (
	"context"
	"encoding/json"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	w *kafka.Writer
}

func (p *Publisher) PublishOrderEvent(ctx context.Context, event *models.OrderEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.Publish(ctx, b)
}


