package ordersackconsumer

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

type ackProcessor interface {
	Handle(ctx context.Context, ack *models.OrderAck) error
}

type OrdersAckConsumer struct {
	processor   ackProcessor
	kafkaBroker []string
	topicName   string
	groupID     string
}

func NewOrdersAckConsumer(processor ackProcessor, kafkaBroker []string, topicName, groupID string) *OrdersAckConsumer {
	return &OrdersAckConsumer{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
		groupID:     groupID,
	}
}


