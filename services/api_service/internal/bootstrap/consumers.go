package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/services/api_service/config"
	ordersackconsumer "github.com/ivkhr/orderstream/services/api_service/internal/consumer/orders_ack_consumer"
	ordersackprocessor "github.com/ivkhr/orderstream/services/api_service/internal/services/processors/orders_ack_processor"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/kafkastorage"
)

func InitAckConsumer(cfg *config.Config, processor *ordersackprocessor.Processor) *ordersackconsumer.Consumer {
	return ordersackconsumer.New(
		processor,
		kafkastorage.NewKafkaGoReaderFactory(),
		strings.Split(cfg.KafkaBrokers, ","),
		cfg.OrdersAckTopic,
		"api-ack-group",
	)
}
