package orderseventprocessor

import (
	"context"
	"log/slog"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

func (p *Processor) Handle(ctx context.Context, event *models.OrderEvent) error {
	slog.Info("worker.orders_event_processor.handle", "order_id", event.OrderID, "status", event.Status)
	res, err := p.svc.HandleOrderEvent(ctx, event)
	if err != nil {
		return err
	}
	if err := p.results.SetOrderResult(ctx, res); err != nil {
		return err
	}
	slog.Info("worker.orders_event_processor.result_saved", "order_id", res.OrderID, "status", res.Status)
	return p.ackPub.PublishOrderAck(ctx, &models.OrderAck{
		OrderID:     res.OrderID,
		Status:      res.Status,
		ProcessedAt: res.ProcessedAt,
	})
}
