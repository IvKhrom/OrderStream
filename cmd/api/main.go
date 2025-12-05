package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/adapters/kafka"
	"github.com/ivkhr/orderstream/internal/adapters/postgres"
	"github.com/ivkhr/orderstream/internal/domain"
)

type createReq struct {
	ExternalID string          `json:"external_id"`
	UserID     string          `json:"user_id"`
	Payload    json.RawMessage `json:"payload"`
}

type ackWaiter struct {
	ch      chan bool
	orderID string
}

var (
	ackWaiters = make(map[string]chan bool)
	ackMutex   sync.RWMutex
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
	producer := kafka.NewProducer([]string{brokers}, "orders.events")
	defer producer.Close()

	// Start ack listener
	go startAckListener(brokers)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		var req createReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.ExternalID == "" || req.UserID == "" {
			http.Error(w, "external_id and user_id are required", http.StatusBadRequest)
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		// Create internal order ID and bucket
		oid := uuid.New()
		bucket := domain.BucketFromUUID(oid, 4)

		ord := &domain.Order{
			OrderID:    oid,
			ExternalID: req.ExternalID, // Используется при сохранении в БД
			UserID:     userID,         // Используется при сохранении в БД
			Amount:     0.0,
			Payload:    req.Payload,
			Status:     "new",
			Bucket:     bucket,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Save to database first
		if err := db.CreateOrder(r.Context(), ord); err != nil {
			log.Error().Err(err).Msg("create order db")
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		// Create ack waiter
		ackCh := make(chan bool, 1)
		ackMutex.Lock()
		ackWaiters[oid.String()] = ackCh
		ackMutex.Unlock()

		defer func() {
			ackMutex.Lock()
			delete(ackWaiters, oid.String())
			ackMutex.Unlock()
		}()

		// Create order event
		event := domain.OrderEvent{
			EventID:    "0", // 0 = create operation
			OrderID:    oid.String(),
			ExternalID: req.ExternalID, // Используется в событии
			UserID:     req.UserID,
			Payload:    req.Payload,
			Status:     "new",
			Timestamp:  time.Now().UTC(),
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			http.Error(w, "event marshal error", http.StatusInternalServerError)
			return
		}

		// Publish to Kafka
		if err := producer.Publish(r.Context(), eventBytes); err != nil {
			log.Error().Err(err).Msg("kafka publish")
			http.Error(w, "kafka error", http.StatusInternalServerError)
			return
		}

		// Wait for ack with timeout
		select {
		case <-ackCh:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"order_id":    ord.OrderID.String(),
				"external_id": ord.ExternalID, // Используется в ответе
				"status":      "created",
			})
		case <-time.After(30 * time.Second):
			http.Error(w, "processing timeout", http.StatusGatewayTimeout)
		case <-r.Context().Done():
			return
		}
	})

	r.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		oid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		ord, err := db.GetByID(r.Context(), oid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ord)
	})

	r.Get("/orders/by-external/{external}", func(w http.ResponseWriter, r *http.Request) {
		external := chi.URLParam(r, "external")
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		ord, err := db.GetByExternalID(r.Context(), external, userID)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ord)
	})

	// Новые ручки для управления заказами
	r.Put("/orders/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		// Validate status
		validStatuses := map[string]bool{
			"processing": true,
			"done":       true,
			"cancelled":  true,
			"failed":     true,
		}
		if !validStatuses[req.Status] {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}

		// Create update event
		event := domain.OrderEvent{
			EventID:   id, // Используем order_id как event_id для update
			OrderID:   id,
			Status:    req.Status,
			Timestamp: time.Now().UTC(),
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			http.Error(w, "event marshal error", http.StatusInternalServerError)
			return
		}

		// Publish to Kafka
		if err := producer.Publish(r.Context(), eventBytes); err != nil {
			log.Error().Err(err).Msg("kafka publish")
			http.Error(w, "kafka error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"order_id": id,
			"status":   "update_requested",
		})
	})

	r.Delete("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// Create cancellation event
		event := domain.OrderEvent{
			EventID:   id, // Используем order_id как event_id для update
			OrderID:   id,
			Status:    "cancelled",
			Timestamp: time.Now().UTC(),
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			http.Error(w, "event marshal error", http.StatusInternalServerError)
			return
		}

		// Publish to Kafka
		if err := producer.Publish(r.Context(), eventBytes); err != nil {
			log.Error().Err(err).Msg("kafka publish")
			http.Error(w, "kafka error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"order_id": id,
			"status":   "cancellation_requested",
		})
	})

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

func startAckListener(brokers string) {
	consumer := kafka.NewConsumer([]string{brokers}, "orders.ack", "api-ack-group")
	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("ack listener error")
			time.Sleep(time.Second)
			continue
		}

		var ack domain.OrderAck
		if err := json.Unmarshal(msg.Value, &ack); err != nil {
			log.Error().Err(err).Msg("ack unmarshal error")
			continue
		}

		ackMutex.RLock()
		ch, exists := ackWaiters[ack.OrderID]
		ackMutex.RUnlock()

		if exists {
			select {
			case ch <- true:
				log.Info().Str("order_id", ack.OrderID).Msg("ack received")
			default:
			}
		}
	}
}
