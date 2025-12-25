package orders

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/shared/models"
)

var (
	ErrDeletedConflict = errors.New("order already deleted")
	ErrNotFound        = errors.New("order not found")
	ErrResultNotReady  = errors.New("order result not found in redis")
)

type Storage interface {
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, o *models.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	DeleteOrder(ctx context.Context, id string) error
}

type EventsPublisher interface {
	PublishOrderEvent(ctx context.Context, event *models.OrderEvent) error
}

type ResultsStore interface {
	GetOrderAck(ctx context.Context, orderID string) (*models.OrderAck, bool, error)
}

type Service struct {
	storage   Storage
	eventsPub EventsPublisher
	results   ResultsStore

	ackCoordinator AckCoordinator
}

func New(storage Storage, eventsPub EventsPublisher, results ResultsStore, ackReg AckWaitRegistry, ackWaitTimeout time.Duration) *Service {
	return &Service{
		storage:        storage,
		eventsPub:      eventsPub,
		results:        results,
		ackCoordinator: NewAckCoordinator(ackReg, ackWaitTimeout),
	}
}

// Upsert принимает заказ, сохраняет сырой заказ в API БД, публикует событие в Kafka,
// ждёт ACK и возвращает результат из Redis (как в заданном workflow).
func (s *Service) Upsert(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error) {
	if orderID == "" || orderID == "0" {
		oid := uuid.New()
		bucket := models.BucketFromUUID(oid, 4)

		ord := &models.Order{
			OrderID:   oid,
			UserID:    userID,
			Amount:    0,
			Payload:   payload,
			Status:    "new",
			Bucket:    bucket,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.storage.Create(ctx, ord); err != nil {
			return "", "", err
		}

		if err := s.publishAndWait(ctx, ord.OrderID.String(), userID.String(), ord.Payload, ord.Status); err != nil {
			return "", "", err
		}

		ack, ok, err := s.results.GetOrderAck(ctx, ord.OrderID.String())
		if err != nil {
			return ord.OrderID.String(), "", err
		}
		if !ok {
			return ord.OrderID.String(), "", ErrResultNotReady
		}
		return ord.OrderID.String(), ack.Status, nil
	}

	targetID, err := uuid.Parse(orderID)
	if err != nil {
		return "", "", err
	}

	existing, err := s.storage.GetByID(ctx, targetID)
	if err != nil {
		return "", "", ErrNotFound
	}
	if existing.Status == "deleted" {
		return existing.OrderID.String(), "deleted", ErrDeletedConflict
	}

	if status == "deleted" {
		if err := s.storage.DeleteOrder(ctx, existing.OrderID.String()); err != nil {
			return "", "", err
		}
		if err := s.publishAndWait(ctx, existing.OrderID.String(), existing.UserID.String(), existing.Payload, "deleted"); err != nil {
			return "", "", err
		}
		ack, ok, err := s.results.GetOrderAck(ctx, existing.OrderID.String())
		if err != nil {
			return existing.OrderID.String(), "", err
		}
		if !ok {
			return existing.OrderID.String(), "", ErrResultNotReady
		}
		return existing.OrderID.String(), ack.Status, nil
	}

	existing.Payload = payload
	existing.UpdatedAt = time.Now()
	if err := s.storage.Update(ctx, existing); err != nil {
		return "", "", err
	}
	if err := s.publishAndWait(ctx, existing.OrderID.String(), existing.UserID.String(), existing.Payload, existing.Status); err != nil {
		return "", "", err
	}

	ack, ok, err := s.results.GetOrderAck(ctx, existing.OrderID.String())
	if err != nil {
		return existing.OrderID.String(), "", err
	}
	if !ok {
		return existing.OrderID.String(), "", ErrResultNotReady
	}
	return existing.OrderID.String(), ack.Status, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *Service) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	return s.storage.GetByExternalID(ctx, externalID, userID)
}

func (s *Service) publishAndWait(ctx context.Context, orderID, userID string, payload json.RawMessage, status string) error {
	evExternal := payloadID(payload)
	event := models.OrderEvent{
		OrderID:    orderID,
		ExternalID: evExternal,
		UserID:     userID,
		Payload:    payload,
		Status:     status,
		Timestamp:  time.Now().UTC(),
	}

	publish := func() error {
		return s.eventsPub.PublishOrderEvent(ctx, &event)
	}
	if s.ackCoordinator == nil {
		return publish()
	}
	return s.ackCoordinator.ExecuteAndWait(ctx, orderID, publish)
}


