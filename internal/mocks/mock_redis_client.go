package mocks

import (
	"context"
	"time"

	rdb "github.com/go-redis/redis/v8"
)

// Ручной мок RedisClient для тестов адаптера Redis.

type MockStringCmd struct {
	val []byte
	err error
}

func (m *MockStringCmd) Result() (string, error) { return string(m.val), m.err }
func (m *MockStringCmd) Bytes() ([]byte, error)  { return m.val, m.err }

type MockStatusCmd struct {
	err error
}

func (m *MockStatusCmd) Result() (string, error) { return "", m.err }
func (m *MockStatusCmd) Err() error              { return m.err }

type MockRedisClient struct {
	SetFunc   func(ctx context.Context, key string, value interface{}, ttl time.Duration) *rdb.StatusCmd
	GetFunc   func(ctx context.Context, key string) *rdb.StringCmd
	CloseFunc func() error
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *rdb.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, ttl)
	}
	return &rdb.StatusCmd{}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *rdb.StringCmd {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return &rdb.StringCmd{}
}

func (m *MockRedisClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}
