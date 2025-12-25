package resultsredis

import (
	"context"
	"fmt"
	"time"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/redisstorage"
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

func (s *Storage) SetOrderResult(ctx context.Context, res *models.OrderResult) error {
	return s.rs.SetJSON(ctx, key(res.OrderID), res, s.ttl)
}
