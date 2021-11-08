package cron

import (
	"sync"

	"github.com/neghoda/api/src/config"
	"github.com/neghoda/api/src/repo/ssga"
	"github.com/neghoda/api/src/service"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var once = &sync.Once{}

type Wrapper struct {
	cfg      *config.Config
	fundRepo service.FundRepo
	ssgaRepo *ssga.Client
}

func NewWrapper(
	cfg *config.Config,
	fundRepo service.FundRepo,
	ssgaRepo *ssga.Client,
) *Wrapper {
	return &Wrapper{
		cfg:      cfg,
		fundRepo: fundRepo,
		ssgaRepo: ssgaRepo,
	}
}

func (cw *Wrapper) Setup() error {
	var err error

	once.Do(func() {
		c := cron.New()

		_, err = c.AddFunc(cw.cfg.CronConfig, cw.syncFundsData)
		if err != nil {
			log.Errorf("Failed to setup cron job (syncFundsData): %v", err)
		}

		c.Start()
	})

	return err
}
