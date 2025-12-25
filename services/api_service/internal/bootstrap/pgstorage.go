package bootstrap

import (
	"github.com/ivkhr/orderstream/services/api_service/config"
	apiservice "github.com/ivkhr/orderstream/services/api_service/internal/services/api_service"
	"github.com/ivkhr/orderstream/services/api_service/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) (apiservice.Storage, error) {
	return pgstorage.NewPGStorge(cfg.PostgresDSN)
}
