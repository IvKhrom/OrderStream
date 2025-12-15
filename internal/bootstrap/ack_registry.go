package bootstrap

import "github.com/ivkhr/orderstream/internal/services/ordersService"

func InitAckRegistry() *ordersService.AckRegistry {
	return ordersService.NewAckRegistry()
}


