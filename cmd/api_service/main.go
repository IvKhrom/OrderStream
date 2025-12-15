package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/bootstrap"
	"log/slog"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}
	slog.Info("конфигурация", "KAFKA_BROKERS", cfg.KafkaBrokers, "POSTGRES_DSN", cfg.PostgresDSN, "API_PORT", cfg.ApiPort)
	if strings.Contains(cfg.KafkaBrokers, "localhost:9092") {
		slog.Warn("KAFKA_BROKERS=localhost:9092 часто приводит к ошибке lookup kafka из-за advertised listeners; для запуска с хоста используйте localhost:29092")
	}
	if strings.Contains(cfg.PostgresDSN, "localhost:5432") {
		slog.Warn("POSTGRES_DSN указывает на localhost:5432; если Postgres запущен через docker-compose, то порт обычно 5433 (см. docker-compose.yml)")
	}

	pg := bootstrap.InitPGStorage(cfg)
	eventsPublisher := bootstrap.InitOrdersEventsPublisher(cfg)
	ackRegistry := bootstrap.InitAckRegistry()
	ordersService := bootstrap.InitOrdersService(pg, eventsPublisher, ackRegistry, cfg)

	ackProcessor := bootstrap.InitOrdersAckProcessor(ackRegistry)
	ackConsumer := bootstrap.InitOrdersAckConsumer(cfg, ackProcessor)

	api := bootstrap.InitOrdersServiceAPI(ordersService)
	bootstrap.AppRun(*api, ackConsumer, cfg.ApiPort)
}


