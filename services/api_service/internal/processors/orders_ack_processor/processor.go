package ordersackprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/shared/models"
)

type Notifier interface {
	Notify(orderID string)
}

type Processor struct {
	notifier Notifier
}

func New(notifier Notifier) *Processor {
	return &Processor{notifier: notifier}
}

func (p *Processor) Handle(ctx context.Context, ack *models.OrderAck) error {
	_ = ctx
	p.notifier.Notify(ack.OrderID)
	return nil
}


