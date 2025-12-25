package main

import (
	"fmt"

	"github.com/ivkhr/orderstream/services/api_service/config"
	"github.com/ivkhr/orderstream/services/api_service/internal/bootstrap"
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
)

func main() {
	cfg, err := (config.EnvLoader{}).Load()
	if err != nil {
		panic(fmt.Sprintf("config error: %v", err))
	}

	pg, err := bootstrap.InitPGStorage(cfg)
	if err != nil {
		panic(err)
	}

	redis := bootstrap.InitRedis(cfg)
	eventsProducer := bootstrap.InitOrdersEventsProducer(cfg)

	ackRegistry := apiservice.NewAckRegistry()
	results := bootstrap.InitResultsStore(redis)
	ordersService := apiservice.New(pg, eventsProducer, results, ackRegistry, cfg.AckWaitTimeout)

	ackProcessor := bootstrap.InitAckProcessor(ordersService)
	ackConsumer := bootstrap.InitAckConsumer(cfg, ackProcessor)

	apiLayer := bootstrap.InitOrdersServiceAPI(ordersService)

	if err := bootstrap.RunHTTP(cfg.HTTPPort, apiLayer, ackConsumer); err != nil {
		panic(err)
	}
}
