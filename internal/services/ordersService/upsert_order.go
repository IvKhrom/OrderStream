package ordersService

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

var (
	ErrDeletedConflict = errors.New("order already deleted")
	ErrNotFound        = errors.New("order not found")
)

func (s *OrdersService) UpsertOrder(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error) {
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

		if err := s.publishAndMaybeWaitAck(ctx, ord.OrderID.String(), userID.String(), ord.Payload, ord.Status); err != nil {
			return "", "", err
		}
		return ord.OrderID.String(), "created", nil
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
		if err := s.publishAndMaybeWaitAck(ctx, existing.OrderID.String(), existing.UserID.String(), existing.Payload, "deleted"); err != nil {
			return "", "", err
		}
		return existing.OrderID.String(), "deleted", nil
	}

	existing.Payload = payload
	existing.UpdatedAt = time.Now()
	if err := s.storage.Update(ctx, existing); err != nil {
		return "", "", err
	}

	if err := s.publishAndMaybeWaitAck(ctx, existing.OrderID.String(), existing.UserID.String(), existing.Payload, existing.Status); err != nil {
		return "", "", err
	}

	return existing.OrderID.String(), "updated", nil
}

func (s *OrdersService) publishAndMaybeWaitAck(ctx context.Context, orderID, userID string, payload json.RawMessage, status string) error {
	evExternal := payloadID(payload)
	event := models.OrderEvent{
		OrderID:    orderID,
		ExternalID: evExternal,
		UserID:     userID,
		Payload:    payload,
		Status:     status,
		Timestamp:  time.Now().UTC(),
	}
	evb, _ := json.Marshal(event)

	var (
		ch      <-chan struct{}
		cleanup func()
	)
	if s.ackRegistry != nil {
		ch, cleanup = s.ackRegistry.Register(orderID)
		defer cleanup()
	}

	if err := s.eventsPub.Publish(ctx, evb); err != nil {
		return err
	}

	if ch == nil {
		return nil
	}

	select {
	case <-ch:
		return nil
	case <-time.After(s.ackWaitTimeout):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
