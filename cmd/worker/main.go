package main

import (
	"fmt"
	"os"

	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	pg := bootstrap.InitPGStorage(cfg)
	ackProducer := bootstrap.InitOrdersAckPublisher(cfg)

	ordersService := bootstrap.InitOrdersService(pg, nil, nil, cfg)
	eventsProcessor := bootstrap.InitOrdersEventProcessor(ordersService, ackProducer)
	eventsConsumer := bootstrap.InitOrdersEventsConsumer(cfg, eventsProcessor)

	bootstrap.WorkerRun(eventsConsumer)
}
