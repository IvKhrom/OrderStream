package orderseventprocessor

import (
	"context"
	"encoding/json"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *OrdersEventProcessor) Handle(ctx context.Context, event *models.OrderEvent) error {
	ack, err := p.ordersService.HandleOrderEvent(ctx, event)
	if err != nil {
		return err
	}
	ackBytes, _ := json.Marshal(ack)
	return p.ackPublisher.Publish(ctx, ackBytes)
}


