package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ivkhr/orderstream/internal/adapters/kafka"
	"github.com/ivkhr/orderstream/internal/adapters/postgres"
	"github.com/ivkhr/orderstream/internal/domain"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:upvel123@localhost:5433/orderstream?sslmode=disable"
	}
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("connect pg")
	}
	defer db.Close(context.Background())

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	groupID := os.Getenv("WORKER_GROUP")
	if groupID == "" {
		groupID = "order-workers"
	}

	consumer := kafka.NewConsumer([]string{brokers}, "orders.events", groupID)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	for {
		m, err := consumer.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info().Msg("shutting down consumer")
				break
			}
			log.Error().Err(err).Msg("read msg")
			time.Sleep(time.Second)
			continue
		}

		log.Info().Msgf("got msg: %s", string(m.Value))

		// Parse event
		var event domain.OrderEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Error().Err(err).Msg("unmarshal event")
			continue
		}

		// Process based on event type
		if event.EventID == "0" {
			// Create operation
			if err := processCreateEvent(ctx, db, brokers, event); err != nil {
				log.Error().Err(err).Msg("process create event")
			}
		} else {
			// Update operation
			if err := processUpdateEvent(ctx, db, brokers, event); err != nil {
				log.Error().Err(err).Msg("process update event")
			}
		}
	}
}

func processCreateEvent(ctx context.Context, db *postgres.Postgres, brokers string, event domain.OrderEvent) error {
	log.Info().Str("order_id", event.OrderID).Msg("processing create event")

	// Update status to processing
	if err := db.UpdateStatus(ctx, event.OrderID, "processing"); err != nil {
		return err
	}

	// Simulate work (in real app: validate, calculate amount, etc.)
	time.Sleep(2 * time.Second)

	// Update status to done
	if err := db.UpdateStatus(ctx, event.OrderID, "done"); err != nil {
		return err
	}

	// Send ack
	ack := domain.OrderAck{
		EventID:     event.EventID,
		OrderID:     event.OrderID,
		Status:      "processed",
		ProcessedAt: time.Now().UTC(),
	}

	ackBytes, _ := json.Marshal(ack)
	if err := kafka.PublishRaw(ctx, []string{brokers}, "orders.ack", ackBytes); err != nil {
		log.Error().Err(err).Msg("publish ack")
		return err
	}

	log.Info().Str("order_id", event.OrderID).Msg("create event processed")
	return nil
}

func processUpdateEvent(ctx context.Context, db *postgres.Postgres, brokers string, event domain.OrderEvent) error {
	log.Info().Str("order_id", event.OrderID).Str("status", event.Status).Msg("processing update event")

	// Update status in database
	if err := db.UpdateStatus(ctx, event.OrderID, event.Status); err != nil {
		return err
	}

	// For cancellation, we might do additional cleanup
	if event.Status == "cancelled" {
		// Additional cancellation logic here
		log.Info().Str("order_id", event.OrderID).Msg("order cancelled")
	}

	// Send ack for update
	ack := domain.OrderAck{
		EventID:     event.EventID,
		OrderID:     event.OrderID,
		Status:      "updated",
		ProcessedAt: time.Now().UTC(),
	}

	ackBytes, _ := json.Marshal(ack)
	if err := kafka.PublishRaw(ctx, []string{brokers}, "orders.ack", ackBytes); err != nil {
		log.Error().Err(err).Msg("publish ack")
		return err
	}

	log.Info().Str("order_id", event.OrderID).Msg("update event processed")
	return nil
}
