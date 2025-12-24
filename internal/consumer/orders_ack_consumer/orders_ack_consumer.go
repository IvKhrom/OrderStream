package ordersackconsumer

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

type ackProcessor interface {
	Handle(ctx context.Context, ack *models.OrderAck) error
}

type OrdersAckConsumer struct {
	processor   ackProcessor
	readerFactory kafkastorage.ReaderFactory
	kafkaBroker []string
	topicName   string
	groupID     string
}

func NewOrdersAckConsumer(processor ackProcessor, readerFactory kafkastorage.ReaderFactory, kafkaBroker []string, topicName, groupID string) *OrdersAckConsumer {
	return &OrdersAckConsumer{
		processor:   processor,
		readerFactory: readerFactory,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
		groupID:     groupID,
	}
}


