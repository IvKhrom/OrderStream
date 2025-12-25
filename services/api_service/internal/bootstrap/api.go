package bootstrap

import (
	ordersapi "github.com/ivkhr/orderstream/services/api_service/internal/api/orders_service_api"
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
)

func InitOrdersServiceAPI(ordersService *apiservice.Service) *ordersapi.OrdersServiceAPI {
	return ordersapi.NewOrdersServiceAPI(ordersService)
}
