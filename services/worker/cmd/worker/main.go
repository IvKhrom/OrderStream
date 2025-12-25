package main

import (
	"fmt"
	"time"

	"github.com/ivkhr/orderstream/services/worker/config"
	"github.com/ivkhr/orderstream/services/worker/internal/bootstrap"
	workersvc "github.com/ivkhr/orderstream/services/worker/internal/services/worker"
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

	redis := bootstrap.InitRedis(cfg.RedisAddr)
	results := bootstrap.InitResultWriter(redis, 10*time.Minute)
	ackPublisher := bootstrap.InitAckPublisher(cfg)

	svc := workersvc.New(pg)
	proc := bootstrap.InitOrdersEventProcessor(svc, results, ackPublisher)
	cons := bootstrap.InitOrdersEventsConsumer(cfg, proc)

	bootstrap.Run(cons)
}
