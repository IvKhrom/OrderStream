package ordersackconsumer

import (
	"context"

	"github.com/ivkhr/orderstream/shared/models"
	"github.com/ivkhr/orderstream/shared/storage/kafkastorage"
)

type Processor interface {
	Handle(ctx context.Context, ack *models.OrderAck) error
}

type Consumer struct {
	processor     Processor
	readerFactory kafkastorage.ReaderFactory
	brokers       []string
	topic         string
	groupID       string
}

func New(processor Processor, rf kafkastorage.ReaderFactory, brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		processor:     processor,
		readerFactory: rf,
		brokers:       brokers,
		topic:         topic,
		groupID:       groupID,
	}
}


