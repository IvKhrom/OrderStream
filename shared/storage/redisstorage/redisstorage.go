package redisstorage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client — минимальный интерфейс для redis.Client (чтобы можно было мокать).
type Client interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type Storage struct {
	c Client
}

func New(addr string) *Storage {
	return &Storage{
		c: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func NewWithClient(c Client) *Storage {
	return &Storage{c: c}
}

func (s *Storage) SetJSON(ctx context.Context, key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.c.Set(ctx, key, b, ttl).Err()
}

func (s *Storage) GetJSON(ctx context.Context, key string, out any) (bool, error) {
	val, err := s.c.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(val, out); err != nil {
		return false, err
	}
	return true, nil
}


