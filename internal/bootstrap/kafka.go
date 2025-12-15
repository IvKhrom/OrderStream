package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

func InitOrdersEventsPublisher(cfg *config.Config) *kafkastorage.Publisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic)
}

func InitOrdersAckPublisher(cfg *config.Config) *kafkastorage.Publisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersAckTopic)
}


