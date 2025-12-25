package api_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/services/api_service/internal/models"
)

func (s *Service) UpsertOrder(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error) {
	return s.Upsert(ctx, orderID, userID, status, payload)
}

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


