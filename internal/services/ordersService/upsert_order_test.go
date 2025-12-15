package ordersService

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

type memStorage struct {
	mu     sync.Mutex
	orders map[uuid.UUID]*models.Order
}

func newMemStorage() *memStorage {
	return &memStorage{orders: make(map[uuid.UUID]*models.Order)}
}

func (m *memStorage) Create(ctx context.Context, o *models.Order) error {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders[o.OrderID] = o
	return nil
}
func (m *memStorage) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	if o, ok := m.orders[id]; ok {
		clone := *o
		return &clone, nil
	}
	return nil, ErrNotFound
}
func (m *memStorage) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	_ = ctx
	_ = externalID
	_ = userID
	return nil, ErrNotFound
}
func (m *memStorage) Update(ctx context.Context, o *models.Order) error {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders[o.OrderID] = o
	return nil
}
func (m *memStorage) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = ctx
	_ = id
	_ = status
	return nil
}
func (m *memStorage) DeleteOrder(ctx context.Context, id string) error {
	_ = ctx
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if o, ok := m.orders[uid]; ok {
		o.Status = "deleted"
	}
	return nil
}

type stubPublisher struct {
	onPublish func(value []byte)
	retErr    error
}

func (p *stubPublisher) Publish(ctx context.Context, value []byte) error {
	_ = ctx
	if p.onPublish != nil {
		p.onPublish(value)
	}
	return p.retErr
}

func TestOrdersService_UpsertOrder_Create_WaitsAck(t *testing.T) {
	st := newMemStorage()
	reg := NewAckRegistry()

	pub := &stubPublisher{
		onPublish: func(value []byte) {
			var ev models.OrderEvent
			_ = json.Unmarshal(value, &ev)
			reg.Notify(ev.OrderID)
		},
	}

	svc := NewOrdersService(st, pub, reg, 2*time.Second)

	userID := uuid.New()
	payload := json.RawMessage(`{"id":"ext-1","items":[1]}`)

	orderID, status, err := svc.UpsertOrder(context.Background(), "0", userID, "", payload)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if orderID == "" {
		t.Fatalf("expected orderID")
	}
	if status != "created" {
		t.Fatalf("expected status=created, got %q", status)
	}
}

func TestOrdersService_UpsertOrder_Update(t *testing.T) {
	st := newMemStorage()
	svc := NewOrdersService(st, &stubPublisher{}, nil, 0)

	userID := uuid.New()
	oid := uuid.New()
	_ = st.Create(context.Background(), &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Status:    "new",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	newPayload := json.RawMessage(`{"id":"ext","x":1}`)
	gotID, gotStatus, err := svc.UpsertOrder(context.Background(), oid.String(), userID, "", newPayload)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if gotID != oid.String() || gotStatus != "updated" {
		t.Fatalf("ожидали updated для %s, получили %s/%s", oid.String(), gotID, gotStatus)
	}
}

func TestOrdersService_UpsertOrder_Delete(t *testing.T) {
	st := newMemStorage()
	svc := NewOrdersService(st, &stubPublisher{}, nil, 0)

	userID := uuid.New()
	oid := uuid.New()
	_ = st.Create(context.Background(), &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Status:    "new",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	gotID, gotStatus, err := svc.UpsertOrder(context.Background(), oid.String(), userID, "deleted", json.RawMessage(`{"id":"ext"}`))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if gotID != oid.String() || gotStatus != "deleted" {
		t.Fatalf("ожидали deleted для %s, получили %s/%s", oid.String(), gotID, gotStatus)
	}
}

func TestOrdersService_UpsertOrder_DeletedConflict(t *testing.T) {
	st := newMemStorage()
	svc := NewOrdersService(st, &stubPublisher{}, nil, 0)

	userID := uuid.New()
	oid := uuid.New()
	_ = st.Create(context.Background(), &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Status:    "deleted",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_, _, err := svc.UpsertOrder(context.Background(), oid.String(), userID, "", json.RawMessage(`{}`))
	if err != ErrDeletedConflict {
		t.Fatalf("ожидали ErrDeletedConflict, получили %v", err)
	}
}

func TestOrdersService_UpsertOrder_PublishError(t *testing.T) {
	st := newMemStorage()
	pub := &stubPublisher{retErr: errors.New("kafka error")}
	svc := NewOrdersService(st, pub, nil, 0)

	userID := uuid.New()
	_, _, err := svc.UpsertOrder(context.Background(), "0", userID, "", json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("ожидали ошибку публикации")
	}
}

func TestOrdersService_UpsertOrder_InvalidOrderID(t *testing.T) {
	st := newMemStorage()
	svc := NewOrdersService(st, &stubPublisher{}, nil, 0)

	_, _, err := svc.UpsertOrder(context.Background(), "bad-uuid", uuid.New(), "", json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("ожидали ошибку парсинга uuid")
	}
}

func TestOrdersService_UpsertOrder_WaitAckTimeout(t *testing.T) {
	st := newMemStorage()
	reg := NewAckRegistry()
	// publish успешный, но ACK никогда не придёт
	pub := &stubPublisher{}
	svc := NewOrdersService(st, pub, reg, 1*time.Millisecond)

	_, _, err := svc.UpsertOrder(context.Background(), "0", uuid.New(), "", json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("ожидали таймаут ожидания ACK")
	}
}


