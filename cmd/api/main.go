package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ivkhr/orderstream/internal/adapters/kafka"
	"github.com/ivkhr/orderstream/internal/adapters/postgres"
	"github.com/ivkhr/orderstream/internal/api"
	"github.com/ivkhr/orderstream/internal/config"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	db, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("connect pg")
	}
	defer db.Close(context.Background())

	brokers := cfg.KafkaBrokers
	producer := kafka.NewProducer([]string{brokers}, "orders.events")
	defer producer.Close()

	ackConsumer := kafka.NewConsumer([]string{brokers}, "orders.ack", "api-ack-group")

	r := api.NewRouter(db, producer, ackConsumer)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Info().Msgf("api listening %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("serve")
	}
}
