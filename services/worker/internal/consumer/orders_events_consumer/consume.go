package orderseventsconsumer

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
			slog.Error("OrdersEventsConsumer.read error", "error", err.Error())
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


