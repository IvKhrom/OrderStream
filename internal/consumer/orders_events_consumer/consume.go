package orderseventsconsumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *OrdersEventsConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           c.groupID,
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("OrdersEventsConsumer.consume error", "error", err.Error())
			continue
		}

		var event *models.OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			slog.Error("OrdersEventsConsumer.unmarshal error", "error", err.Error())
			continue
		}
		if err := c.processor.Handle(ctx, event); err != nil {
			slog.Error("OrdersEventsConsumer.handle error", "error", err.Error())
		}
	}
}


