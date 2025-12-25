package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/services/worker/config"
	orderseventsconsumer "github.com/ivkhr/orderstream/services/worker/internal/consumer/orders_events_consumer"
	orderseventprocessor "github.com/ivkhr/orderstream/services/worker/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/kafkastorage"
)

func InitOrdersEventsConsumer(cfg *config.Config, processor *orderseventprocessor.Processor) *orderseventsconsumer.Consumer {
	return orderseventsconsumer.New(
		processor,
		kafkastorage.NewKafkaGoReaderFactory(),
		strings.Split(cfg.KafkaBrokers, ","),
		cfg.OrdersEventsTopic,
		cfg.WorkerGroup,
	)
}
