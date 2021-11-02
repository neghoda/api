package service

import (
	"sync"

	"github.com/erp/api/src/config"
	"github.com/erp/api/src/storage/postgres"
	"github.com/erp/api/src/storage/redis"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	redis    *redis.Client
	authRepo *postgres.Connector
	userRepo *postgres.Connector
}

func New(
	cfg *config.Config,
	rds *redis.Client,
	aur *postgres.Connector,
	pr *postgres.Connector,
) *Service {
	once.Do(func() {
		service = &Service{
			redis:    rds,
			cfg:      cfg,
			authRepo: aur,
			userRepo: pr,
		}
	})

	return service
}

func Get() *Service {
	return service
}
