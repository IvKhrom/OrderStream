package resultsredis

import (
	"context"
	"fmt"

	"github.com/ivkhr/orderstream/services/api_service/internal/models"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/redisstorage"
)

type Storage struct {
	rs *redisstorage.Storage
}

func New(rs *redisstorage.Storage) *Storage {
	return &Storage{rs: rs}
}

func key(orderID string) string {
	return fmt.Sprintf("order_result:%s", orderID)
}

func (s *Storage) GetOrderAck(ctx context.Context, orderID string) (*models.OrderAck, bool, error) {
	var ack models.OrderAck
	ok, err := s.rs.GetJSON(ctx, key(orderID), &ack)
	if err != nil || !ok {
		return nil, ok, err
	}
	return &ack, true, nil
}

func (s *Storage) GetOrderResult(ctx context.Context, orderID string) (*models.OrderResult, bool, error) {
	var res models.OrderResult
	ok, err := s.rs.GetJSON(ctx, key(orderID), &res)
	if err != nil || !ok {
		return nil, ok, err
	}
	return &res, true, nil
}


