package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/adapters/kafka"
	"github.com/ivkhr/orderstream/internal/domain"
	"github.com/ivkhr/orderstream/internal/repository"
)

// вспомогательная функция: пытается извлечь идентификатор маркетплейса из JSON payload (поле "id")
func payloadID(payload json.RawMessage) string {
	var tmp map[string]any
	if err := json.Unmarshal(payload, &tmp); err != nil {
		return ""
	}
	if v, ok := tmp["id"]; ok {
		switch t := v.(type) {
		case string:
			return t
		case float64:
			// числовой id — форматируем без дробной части, если это целое число
			return fmt.Sprintf("%v", t)
		default:
			return ""
		}
	}
	return ""
}

// NewRouter создаёт HTTP-маршрутизатор с внедрёнными зависимостями.
// Параметр ackConsumer может быть nil — тогда ожидание подтверждений отключено.
func NewRouter(repo repository.OrderRepository, producer kafka.ProducerClient, ackConsumer kafka.ConsumerClient) http.Handler {
	r := chi.NewRouter()

	var (
		ackWaiters = make(map[string]chan bool)
		ackMutex   sync.RWMutex
	)

	// Запускаем слушатель подтверждений (если он передан)
	if ackConsumer != nil {
		go func() {
			defer ackConsumer.Close()
			for {
				msg, err := ackConsumer.ReadMessage(context.Background())
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				var ack domain.OrderAck
				if err := json.Unmarshal(msg.Value, &ack); err != nil {
					continue
				}
				ackMutex.RLock()
				ch, ok := ackWaiters[ack.OrderID]
				ackMutex.RUnlock()
				if ok {
					select {
					case ch <- true:
					default:
					}
				}
			}
		}()
	}

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Post("/orders", func(w http.ResponseWriter, req *http.Request) {
		var in struct {
			OrderID string          `json:"order_id"` // пустой или "0" -> создать, иначе считать как UUID для обновления
			UserID  string          `json:"user_id"`
			Status  string          `json:"status,omitempty"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if in.UserID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		userID, err := uuid.Parse(in.UserID)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		// Создание заказа
		if in.OrderID == "" || in.OrderID == "0" {
			oid := uuid.New()
			bucket := domain.BucketFromUUID(oid, 4)

			ord := &domain.Order{
				OrderID:   oid,
				UserID:    userID,
				Amount:    0,
				Payload:   in.Payload,
				Status:    "new",
				Bucket:    bucket,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := repo.Create(req.Context(), ord); err != nil {
				http.Error(w, "db error", http.StatusInternalServerError)
				return
			}

			// публикуем событие (create) — добавляем внешний id из payload, если он есть
			evExternal := payloadID(ord.Payload)
			event := domain.OrderEvent{
				OrderID:    ord.OrderID.String(),
				ExternalID: evExternal,
				UserID:     ord.UserID.String(),
				Payload:    ord.Payload,
				Status:     ord.Status,
				Timestamp:  time.Now().UTC(),
			}
			evb, _ := json.Marshal(event)
			if err := producer.Publish(req.Context(), evb); err != nil {
				http.Error(w, "kafka error", http.StatusInternalServerError)
				return
			}

			// ожидаем подтверждение (ACK), если подключён ackConsumer
			if ackConsumer != nil {
				ch := make(chan bool, 1)
				ackMutex.Lock()
				ackWaiters[ord.OrderID.String()] = ch
				ackMutex.Unlock()
				defer func() {
					ackMutex.Lock()
					delete(ackWaiters, ord.OrderID.String())
					ackMutex.Unlock()
				}()

				select {
				case <-ch:
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]string{"order_id": ord.OrderID.String(), "status": "created"})
					return
				case <-time.After(30 * time.Second):
					http.Error(w, "processing timeout", http.StatusGatewayTimeout)
					return
				case <-req.Context().Done():
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"order_id": ord.OrderID.String(), "status": "created"})
			return
		}

		// Путь обновления: если указан OrderID, считаем его внутренним UUID заказа
		targetID, err := uuid.Parse(in.OrderID)
		if err != nil {
			http.Error(w, "invalid order_id for update", http.StatusBadRequest)
			return
		}

		existing, err := repo.GetByID(req.Context(), targetID)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if existing.Status == "deleted" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
			return
		}

		if in.Status == "deleted" {
			if err := repo.DeleteOrder(req.Context(), existing.OrderID.String()); err != nil {
				http.Error(w, "db error", http.StatusInternalServerError)
				return
			}
			// публикуем событие удаления
			evExternal := payloadID(existing.Payload)
			event := domain.OrderEvent{
				OrderID:    existing.OrderID.String(),
				ExternalID: evExternal,
				UserID:     existing.UserID.String(),
				Payload:    existing.Payload,
				Status:     "deleted",
				Timestamp:  time.Now().UTC(),
			}
			evb, _ := json.Marshal(event)
			_ = producer.Publish(req.Context(), evb)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
			return
		}

		// применяем обновление (заменяем payload; поле amount здесь не меняется)
		existing.Payload = in.Payload
		existing.UpdatedAt = time.Now()

		if err := repo.Update(req.Context(), existing); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		// публикуем событие обновления
		evExternal := payloadID(existing.Payload)
		event := domain.OrderEvent{
			OrderID:    existing.OrderID.String(),
			ExternalID: evExternal,
			UserID:     existing.UserID.String(),
			Payload:    existing.Payload,
			Status:     existing.Status,
			Timestamp:  time.Now().UTC(),
		}
		evb, _ := json.Marshal(event)
		if err := producer.Publish(req.Context(), evb); err != nil {
			http.Error(w, "kafka error", http.StatusInternalServerError)
			return
		}

		// ожидаем подтверждение (ACK) по order_id, если доступен ackConsumer
		if ackConsumer != nil {
			ch := make(chan bool, 1)
			ackMutex.Lock()
			ackWaiters[existing.OrderID.String()] = ch
			ackMutex.Unlock()
			defer func() {
				ackMutex.Lock()
				delete(ackWaiters, existing.OrderID.String())
				ackMutex.Unlock()
			}()

			select {
			case <-ch:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"order_id": existing.OrderID.String(), "status": "updated"})
				return
			case <-time.After(30 * time.Second):
				http.Error(w, "processing timeout", http.StatusGatewayTimeout)
				return
			case <-req.Context().Done():
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"order_id": existing.OrderID.String(), "status": "updated"})
	})

	r.Get("/orders/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		oid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		o, err := repo.GetByID(req.Context(), oid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(o)
	})

	r.Get("/orders/by-external/{external}", func(w http.ResponseWriter, req *http.Request) {
		external := chi.URLParam(req, "external")
		userIDStr := req.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		o, err := repo.GetByExternalID(req.Context(), external, userID)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(o)
	})

	return r
}
