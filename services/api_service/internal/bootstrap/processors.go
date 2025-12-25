package bootstrap

import (
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
	ordersackprocessor "github.com/ivkhr/orderstream/services/api_service/internal/services/processors/orders_ack_processor"
)

func InitAckProcessor(ordersService *apiservice.Service) *ordersackprocessor.Processor {
	return ordersackprocessor.New(ordersService)
}
