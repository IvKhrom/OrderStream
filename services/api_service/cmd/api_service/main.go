package main

import (
	"fmt"
	"strings"

	"github.com/ivkhr/orderstream/shared/storage/kafkastorage"
	"github.com/ivkhr/orderstream/shared/storage/pgstorage"
	"github.com/ivkhr/orderstream/shared/storage/redisstorage"
	"github.com/ivkhr/orderstream/services/api_service/internal/api/orders_http"
	"github.com/ivkhr/orderstream/services/api_service/internal/bootstrap"
	"github.com/ivkhr/orderstream/services/api_service/internal/config"
	ordersackconsumer "github.com/ivkhr/orderstream/services/api_service/internal/consumer/orders_ack_consumer"
	ordersackprocessor "github.com/ivkhr/orderstream/services/api_service/internal/processors/orders_ack_processor"
	"github.com/ivkhr/orderstream/services/api_service/internal/services/orders"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/resultsredis"
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

	eventsPub := kafkastorage.NewPublisher(strings.Split(cfg.KafkaBrokers, ","), cfg.OrdersEventsTopic)

	rs := redisstorage.New(cfg.RedisAddr)
	results := resultsredis.New(rs)

	ackReg := orders.NewAckRegistry()
	svc := orders.New(pg, eventsPub, results, ackReg, cfg.AckWaitTimeout)

	ackProc := ordersackprocessor.New(ackReg)
	ackCons := ordersackconsumer.New(
		ackProc,
		kafkastorage.NewKafkaGoReaderFactory(),
		strings.Split(cfg.KafkaBrokers, ","),
		cfg.OrdersAckTopic,
		"api-ack-group",
	)

	handler := ordershttp.New(svc)
	if err := bootstrap.RunHTTP(cfg.HTTPPort, handler, ackCons); err != nil {
		panic(err)
	}
}


