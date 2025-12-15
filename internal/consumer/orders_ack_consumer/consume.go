package ordersackconsumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *OrdersAckConsumer) Consume(ctx context.Context) {
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
			slog.Error("OrdersAckConsumer.consume error", "error", err.Error())
			continue
		}
		var ack *models.OrderAck
		if err := json.Unmarshal(msg.Value, &ack); err != nil {
			slog.Error("OrdersAckConsumer.unmarshal error", "error", err.Error())
			continue
		}
		if err := c.processor.Handle(ctx, ack); err != nil {
			slog.Error("OrdersAckConsumer.handle error", "error", err.Error())
		}
	}
}


