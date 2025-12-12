package redis

import (
	"context"
	"time"

	rdb "github.com/go-redis/redis/v8"
)

//go:generate mockery --name RedisClient --output ../mocks --outpkg mocks --case underscore

// RedisClient is an interface to allow mocking redis client in tests.
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *rdb.StatusCmd
	Get(ctx context.Context, key string) *rdb.StringCmd
	Close() error
}

type RedisCache struct {
	client RedisClient
}

func New(addr string) *RedisCache {
	opt := &rdb.Options{Addr: addr}
	c := rdb.NewClient(opt)
	return &RedisCache{client: c}
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	s, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return s, nil
}
