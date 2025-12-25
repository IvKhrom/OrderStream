package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/services/api_service/config"
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/kafkastorage"
)

func InitOrdersEventsProducer(cfg *config.Config) apiservice.EventsPublisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic)
}
