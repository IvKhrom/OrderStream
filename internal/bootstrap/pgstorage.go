package bootstrap

import (
	"log"

	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
	"github.com/ivkhr/orderstream/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) ordersService.OrdersStorage {
	storage, err := pgstorage.NewPGStorge(cfg.PostgresDSN)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
	}
	return storage
}
