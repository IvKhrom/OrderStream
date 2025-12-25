package bootstrap

import (
	"github.com/ivkhr/orderstream/services/worker/config"
	workersvc "github.com/ivkhr/orderstream/services/worker/internal/services/worker"
	"github.com/ivkhr/orderstream/services/worker/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) (workersvc.Storage, error) {
	return pgstorage.NewPGStorge(cfg.PostgresDSN)
}


