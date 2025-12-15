package orderseventprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

type ordersService interface {
	HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error)
}

type ackPublisher interface {
	Publish(ctx context.Context, value []byte) error
}

type OrdersEventProcessor struct {
	ordersService ordersService
	ackPublisher  ackPublisher
}

func NewOrdersEventProcessor(ordersService ordersService, ackPublisher ackPublisher) *OrdersEventProcessor {
	return &OrdersEventProcessor{
		ordersService: ordersService,
		ackPublisher:  ackPublisher,
	}
}


