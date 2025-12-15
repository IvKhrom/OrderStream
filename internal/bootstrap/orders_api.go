package bootstrap

import (
	server "github.com/ivkhr/orderstream/internal/api/orders_service_api"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
)

func InitOrdersServiceAPI(ordersService *ordersService.OrdersService) *server.OrdersServiceAPI {
	return server.NewOrdersServiceAPI(ordersService)
}


