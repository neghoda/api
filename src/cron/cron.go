package cron

import (
	"sync"

	"github.com/neghoda/api/src/config"
	"github.com/neghoda/api/src/repo/ssga"
	"github.com/neghoda/api/src/storage/postgres"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var once = &sync.Once{}

type CronWrapper struct {
	cfg      *config.Config
	db       *postgres.Connector
	ssgaRepo *ssga.Client
}

func NewCronWrapper(
	cfg *config.Config,
	db *postgres.Connector,
	ssgaRepo *ssga.Client,
) *CronWrapper {
	return &CronWrapper{
		cfg:      cfg,
		db:       db,
		ssgaRepo: ssgaRepo,
	}
}

func (cw *CronWrapper) Setup() error {
	var err error

	once.Do(func() {
		c := cron.New()

		_, err = c.AddFunc(cw.cfg.CronConfig, cw.syncFundsData)
		if err != nil {
			log.Error("Failed to setup cron job (generateMockMetrics): %v", err)
		}

		c.Start()
	})

	return err
}
