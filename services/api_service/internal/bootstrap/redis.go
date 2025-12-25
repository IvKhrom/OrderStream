package bootstrap

import (
	"github.com/ivkhr/orderstream/services/api_service/config"
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/redisstorage"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/resultsredis"
)

func InitRedis(cfg *config.Config) *redisstorage.Storage {
	return redisstorage.New(cfg.RedisAddr)
}

func InitResultsStore(rs *redisstorage.Storage) apiservice.ResultsStore {
	return resultsredis.New(rs)
}
