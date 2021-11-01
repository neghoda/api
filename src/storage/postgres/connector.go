package postgres

import (
	"context"
	"sync"

	"github.com/erp/api/src/config"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	pgv2 "gitlab.yalantis.com/gophers/pg/v2"
)

// Postgres pg connection.
type Postgres struct {
	*pgv2.Connector
	ctx context.Context
}

func New(ctx context.Context, wg *sync.WaitGroup, mainCfg, replicaCfg *config.Postgres) (*Postgres, error) {
	connector := pgv2.New()

	err := connector.Connect(&pg.Options{
		Addr:         mainCfg.Host + ":" + mainCfg.Port,
		User:         mainCfg.User,
		Password:     mainCfg.Password,
		Database:     mainCfg.Name,
		PoolSize:     mainCfg.PoolSize,
		WriteTimeout: mainCfg.WriteTimeout,
		ReadTimeout:  mainCfg.ReadTimeout,
		MaxRetries:   mainCfg.MaxRetries,
	})
	if err != nil {
		log.WithError(err).Error("cannot connect to db")

		return nil, err
	}

	if mainCfg.EnableLogger {
		connector.SetLogger(&LoggerAdapter{log.StandardLogger()})
	}

	replicaOpt1 := &pg.Options{
		Addr:         replicaCfg.Host + ":" + replicaCfg.Port,
		User:         replicaCfg.User,
		Password:     replicaCfg.Password,
		Database:     replicaCfg.Name,
		PoolSize:     replicaCfg.PoolSize,
		WriteTimeout: replicaCfg.WriteTimeout,
		ReadTimeout:  replicaCfg.ReadTimeout,
		MaxRetries:   replicaCfg.MaxRetries,
	}

	err = connector.ConnectReplicas(replicaOpt1)
	if err != nil {
		log.WithError(err).Error("cannot connect to replica")

		return nil, err
	}

	if replicaCfg.EnableLogger {
		connector.SetLoggerReplica(&LoggerAdapter{log.StandardLogger()})
	}

	p := &Postgres{Connector: connector, ctx: ctx}

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		err := connector.Close()
		if err != nil {
			log.Error("close db connection error:", err.Error())

			return
		}

		log.Info("close db connection")
	}()

	return p, nil
}

// Check checks db connection.
func (p *Postgres) Check() (err error) {
	return p.Connector.Health()
}
