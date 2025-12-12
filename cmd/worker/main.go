package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ivkhr/orderstream/internal/adapters/kafka"
	"github.com/ivkhr/orderstream/internal/adapters/postgres"
	"github.com/ivkhr/orderstream/internal/config"
	"github.com/ivkhr/orderstream/internal/domain"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// загрузка конфигурации
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
	groupID := cfg.WorkerGroup

	consumer := kafka.NewConsumer([]string{brokers}, "orders.events", groupID)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Корректное завершение (graceful shutdown)
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

		log.Info().Msgf("получено сообщение: %s", string(m.Value))

		// Разбор события
		var event domain.OrderEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Error().Err(err).Msg("unmarshal event")
			continue
		}

		// Обработка в зависимости от типа события
		// Обработка: если в событии указан order_id == "0" или пустой — создаём заказ,
		// иначе считаем, что это обновление существующего заказа.
		if event.OrderID == "" || event.OrderID == "0" {
			// Операция создания
			if err := processCreateEvent(ctx, db, brokers, event); err != nil {
				log.Error().Err(err).Msg("process create event")
			}
		} else {
			// Операция обновления
			if err := processUpdateEvent(ctx, db, brokers, event); err != nil {
				log.Error().Err(err).Msg("process update event")
			}
		}
	}
}

func processCreateEvent(ctx context.Context, db *postgres.Postgres, brokers string, event domain.OrderEvent) error {
	log.Info().Str("order_id", event.OrderID).Msg("processing create event")

	// Если order_id в событии равен "0" или пустой, создаём запись в БД и используем новый UUID
	if event.OrderID == "" || event.OrderID == "0" {
		uid := uuid.New()
		// Попытка распарсить user_id; если не удалось, используем пустой UUID
		var userUUID uuid.UUID
		if parsed, err := uuid.Parse(event.UserID); err == nil {
			userUUID = parsed
		}

		bucket := domain.BucketFromUUID(uid, 4)
		ord := &domain.Order{
			OrderID:   uid,
			UserID:    userUUID,
			Amount:    0,
			Payload:   event.Payload,
			Status:    "processing",
			Bucket:    bucket,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.CreateOrder(ctx, ord); err != nil {
			return err
		}

		// Имитируем обработку
		time.Sleep(2 * time.Second)

		// Пометить как done
		if err := db.UpdateStatus(ctx, uid.String(), "done"); err != nil {
			return err
		}

		// Отправляем ACK с новым order_id
		ack := domain.OrderAck{
			OrderID:     uid.String(),
			Status:      "processed",
			ProcessedAt: time.Now().UTC(),
		}
		ackBytes, _ := json.Marshal(ack)
		if err := kafka.PublishRaw(ctx, []string{brokers}, "orders.ack", ackBytes); err != nil {
			log.Error().Err(err).Msg("publish ack")
			return err
		}

		log.Info().Str("order_id", uid.String()).Msg("create event processed")
		return nil
	}

	// Если order_id указан — просто обновляем статус как раньше
	if err := db.UpdateStatus(ctx, event.OrderID, "processing"); err != nil {
		return err
	}

	// Имитируем обработку
	time.Sleep(2 * time.Second)

	if err := db.UpdateStatus(ctx, event.OrderID, "done"); err != nil {
		return err
	}

	ack := domain.OrderAck{
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

	// Обновляем статус в базе данных
	if err := db.UpdateStatus(ctx, event.OrderID, event.Status); err != nil {
		return err
	}

	if event.Status == "cancelled" {
		log.Info().Str("order_id", event.OrderID).Msg("order cancelled")
	}

	// Отправляем подтверждение (ACK) для обновления
	ack := domain.OrderAck{
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
