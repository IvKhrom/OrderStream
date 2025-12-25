package bootstrap

import (
	"time"

	orderseventprocessor "github.com/ivkhr/orderstream/services/worker/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/resultsredis"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/redisstorage"
)

func InitRedis(addr string) *redisstorage.Storage {
	return redisstorage.New(addr)
}

func InitResultWriter(rs *redisstorage.Storage, ttl time.Duration) orderseventprocessor.ResultWriter {
	return resultsredis.New(rs, ttl)
}


