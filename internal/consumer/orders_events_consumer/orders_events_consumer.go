package orderseventsconsumer

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

type eventsProcessor interface {
	Handle(ctx context.Context, event *models.OrderEvent) error
}

type OrdersEventsConsumer struct {
	processor   eventsProcessor
	kafkaBroker []string
	topicName   string
	groupID     string
}

func NewOrdersEventsConsumer(processor eventsProcessor, kafkaBroker []string, topicName, groupID string) *OrdersEventsConsumer {
	return &OrdersEventsConsumer{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
		groupID:     groupID,
	}
}
