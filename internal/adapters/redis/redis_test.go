package redis

import (
	"context"
	"testing"
	"time"

	rdb "github.com/go-redis/redis/v8"
	"github.com/ivkhr/orderstream/internal/mocks"
)

func TestRedisCacheWithMockClient(t *testing.T) {
	// Создаём Мок клиент
	mc := &mocks.MockRedisClient{}
	mc.SetFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) *rdb.StatusCmd {
		return &rdb.StatusCmd{}
	}
	mc.GetFunc = func(ctx context.Context, key string) *rdb.StringCmd {
		sc := rdb.NewStringCmd(ctx)
		sc.SetVal("value")
		return sc
	}

	rc := &RedisCache{client: mc}
	ctx := context.Background()
	if err := rc.Set(ctx, "k", []byte("v"), 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := rc.Get(ctx, "k")
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if string(b) != "value" {
		t.Fatalf("expected value, got %s", string(b))
	}
	if err := rc.Close(); err != nil {
		t.Fatalf("close returned error: %v", err)
	}
}
