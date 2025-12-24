package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
	"github.com/ivkhr/orderstream/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

// InitOrdersEventsPublisher возвращает паблишер событий заказов как интерфейс верхнего слоя.
func InitOrdersEventsPublisher(cfg *config.Config) ordersService.OrdersEventsPublisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic)
}

func InitOrdersAckPublisher(cfg *config.Config) orderseventprocessor.AckPublisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersAckTopic)
}


