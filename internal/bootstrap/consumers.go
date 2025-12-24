package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/config"
	ordersackconsumer "github.com/ivkhr/orderstream/internal/consumer/orders_ack_consumer"
	orderseventsconsumer "github.com/ivkhr/orderstream/internal/consumer/orders_events_consumer"
	ordersackprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_ack_processor"
	orderseventprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

func InitOrdersEventsConsumer(cfg *config.Config, processor *orderseventprocessor.OrdersEventProcessor) *orderseventsconsumer.OrdersEventsConsumer {
	rf := kafkastorage.NewKafkaGoReaderFactory()
	return orderseventsconsumer.NewOrdersEventsConsumer(processor, rf, strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic, cfg.WorkerGroup)
}

func InitOrdersAckConsumer(cfg *config.Config, processor *ordersackprocessor.OrdersAckProcessor) *ordersackconsumer.OrdersAckConsumer {
	rf := kafkastorage.NewKafkaGoReaderFactory()
	return ordersackconsumer.NewOrdersAckConsumer(processor, rf, strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersAckTopic, "api-ack-group")
}
