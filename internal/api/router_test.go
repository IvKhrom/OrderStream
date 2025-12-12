package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

// mockRepo реализует repository.OrderRepository для тестов.
type mockRepo struct {
	created bool
}

func (m *mockRepo) Create(ctx context.Context, o *domain.Order) error {
	m.created = true
	return nil
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) { return nil, nil }
func (m *mockRepo) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*domain.Order, error) {
	return nil, nil
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id string, status string) error { return nil }
func (m *mockRepo) DeleteOrder(ctx context.Context, id string) error                 { return nil }
func (m *mockRepo) Update(ctx context.Context, o *domain.Order) error                { return nil }

// mockProducer сохраняет опубликованные сообщения.
type mockProducer struct {
	last []byte
	ack  *fakeAckConsumer
}

func (m *mockProducer) Publish(ctx context.Context, value []byte) error {
	m.last = append([]byte(nil), value...)
	// если подключён ack consumer, подготовить соответствующий ACK, чтобы роутер был уведомлён
	if m.ack != nil {
		var ev domain.OrderEvent
		if err := json.Unmarshal(value, &ev); err == nil {
			ack := domain.OrderAck{OrderID: ev.OrderID, Status: "processed", ProcessedAt: time.Now().UTC()}
			b, _ := json.Marshal(ack)
			m.ack.msg = kafka.Message{Value: b}
		}
	}
	return nil
}
func (m *mockProducer) Close() error { return nil }

// fakeAckConsumer возвращает одно сообщение-ACK, затем блокируется.
type fakeAckConsumer struct {
	msg  kafka.Message
	done bool
}

func (f *fakeAckConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if f.done {
		// засыпаем до отмены контекста
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}
	f.done = true
	// небольшая задержка, чтобы смоделировать асинхронное подтверждение
	time.Sleep(20 * time.Millisecond)
	return f.msg, nil
}
func (f *fakeAckConsumer) Close() error { return nil }

func TestPostOrderWaitsAck(t *testing.T) {
	repo := &mockRepo{}
	// Подготавливаем fakeAckConsumer и привязываем его к продюсеру, чтобы Publish создавал соответствующий ACK
	fac := &fakeAckConsumer{}
	prod := &mockProducer{ack: fac}

	r := NewRouter(repo, prod, fac)

	body := map[string]interface{}{"order_id": "0", "user_id": uuid.New().String(), "payload": map[string]interface{}{"items": []interface{}{}}}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(jb))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.True(t, repo.created)

}
