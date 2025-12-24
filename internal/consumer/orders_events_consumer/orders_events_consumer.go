package orderseventsconsumer

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

type eventsProcessor interface {
	Handle(ctx context.Context, event *models.OrderEvent) error
}

type OrdersEventsConsumer struct {
	processor   eventsProcessor
	readerFactory kafkastorage.ReaderFactory
	kafkaBroker []string
	topicName   string
	groupID     string
}

func NewOrdersEventsConsumer(processor eventsProcessor, readerFactory kafkastorage.ReaderFactory, kafkaBroker []string, topicName, groupID string) *OrdersEventsConsumer {
	return &OrdersEventsConsumer{
		processor:   processor,
		readerFactory: readerFactory,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
		groupID:     groupID,
	}
}
