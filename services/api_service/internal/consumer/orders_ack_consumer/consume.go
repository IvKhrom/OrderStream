package ordersackconsumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ivkhr/orderstream/shared/models"
)

func (c *Consumer) Consume(ctx context.Context) {
	r := c.readerFactory.NewReader(c.brokers, c.groupID, c.topic)
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("OrdersAckConsumer.read error", "error", err.Error())
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


