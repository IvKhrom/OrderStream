package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ivkhr/orderstream/shared/storage/kafkastorage"
	"github.com/ivkhr/orderstream/shared/storage/pgstorage"
	"github.com/ivkhr/orderstream/shared/storage/redisstorage"
	"github.com/ivkhr/orderstream/services/worker/internal/bootstrap"
	"github.com/ivkhr/orderstream/services/worker/internal/config"
	orderseventsconsumer "github.com/ivkhr/orderstream/services/worker/internal/consumer/orders_events_consumer"
	orderseventprocessor "github.com/ivkhr/orderstream/services/worker/internal/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/services/worker/internal/services/orders"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/resultsredis"
)

func main() {
	cfg, err := (config.EnvLoader{}).Load()
	if err != nil {
		panic(fmt.Sprintf("config error: %v", err))
	}

	pg, err := pgstorage.NewPGStorge(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}

	ackPub := kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersAckTopic)

	rs := redisstorage.New(cfg.RedisAddr)
	results := resultsredis.New(rs, 10*time.Minute)

	svc := orders.New(pg)
	proc := orderseventprocessor.New(svc, results, ackPub)

	cons := orderseventsconsumer.New(proc, kafkastorage.NewKafkaGoReaderFactory(), strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic, cfg.WorkerGroup)
	bootstrap.Run(cons)
}


