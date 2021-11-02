package service

import (
	"sync"

	"github.com/erp/api/src/config"
	"github.com/erp/api/src/storage/postgres"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	authRepo *postgres.Connector
	userRepo *postgres.Connector
}

func New(
	cfg *config.Config,
	aur *postgres.Connector,
	pr *postgres.Connector,
) *Service {
	once.Do(func() {
		service = &Service{
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
