package ordersackprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *OrdersAckProcessor) Handle(ctx context.Context, ack *models.OrderAck) error {
	_ = ctx
	p.notifier.Notify(ack.OrderID)
	return nil
}


