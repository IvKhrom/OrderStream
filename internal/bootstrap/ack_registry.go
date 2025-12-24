package bootstrap

import "github.com/ivkhr/orderstream/internal/services/ordersService"

func InitAckRegistry() ordersService.AckRegistryContract {
	return ordersService.NewAckRegistry()
}


