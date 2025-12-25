package resultsredis

import (
	"context"
	"fmt"
	"time"

	"github.com/ivkhr/orderstream/shared/models"
	"github.com/ivkhr/orderstream/shared/storage/redisstorage"
)

type Storage struct {
	rs  *redisstorage.Storage
	ttl time.Duration
}

func New(rs *redisstorage.Storage, ttl time.Duration) *Storage {
	return &Storage{rs: rs, ttl: ttl}
}

func key(orderID string) string {
	return fmt.Sprintf("order_result:%s", orderID)
}

func (s *Storage) SetOrderAck(ctx context.Context, ack *models.OrderAck) error {
	return s.rs.SetJSON(ctx, key(ack.OrderID), ack, s.ttl)
}


