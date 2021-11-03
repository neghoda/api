package service

import (
	"sync"

	"github.com/neghoda/api/src/config"
	"github.com/neghoda/api/src/repo/ssga"
	"github.com/neghoda/api/src/storage/postgres"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	db       *postgres.Connector
	ssgaRepo *ssga.Client
}

func New(
	cfg *config.Config,
	db *postgres.Connector,
	ssgaRepo *ssga.Client,
) *Service {
	once.Do(func() {
		service = &Service{
			cfg:      cfg,
			db:       db,
			ssgaRepo: ssgaRepo,
		}
	})

	return service
}

func Get() *Service {
	return service
}
