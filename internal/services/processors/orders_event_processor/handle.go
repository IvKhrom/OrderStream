package orderseventprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *OrdersEventProcessor) Handle(ctx context.Context, event *models.OrderEvent) error {
	ack, err := p.ordersService.HandleOrderEvent(ctx, event)
	if err != nil {
		return err
	}
	return p.ackPublisher.PublishOrderAck(ctx, ack)
}


