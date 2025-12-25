package ordersackprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/services/api_service/internal/models"
)

func (p *Processor) Handle(ctx context.Context, ack *models.OrderAck) error {
	_ = ctx
	p.notifier.Notify(ack.OrderID)
	return nil
}


