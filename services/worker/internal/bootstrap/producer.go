package bootstrap

import (
	"strings"

	"github.com/ivkhr/orderstream/services/worker/config"
	orderseventprocessor "github.com/ivkhr/orderstream/services/worker/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/kafkastorage"
)

func InitAckPublisher(cfg *config.Config) orderseventprocessor.AckPublisher {
	return kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersAckTopic)
}
