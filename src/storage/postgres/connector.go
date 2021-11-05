package postgres

import (
	"github.com/go-pg/pg"
	"github.com/neghoda/api/src/config"
)

// Connector pg connection
type Connector struct {
	conn *pg.DB
}

// NewConn init, overwrite connection
func NewConn(mainCfg *config.Postgres) (*Connector, error) {
	var (
		c   Connector
		err error
	)

	c.conn, err = newConn(mainCfg)

	return &c, err
}

// newConn init connection
func newConn(mainCfg *config.Postgres) (*pg.DB, error) {
	conn := pg.Connect(&pg.Options{
		Addr:         mainCfg.Host + ":" + mainCfg.Port,
		User:         mainCfg.User,
		Password:     mainCfg.Password,
		Database:     mainCfg.Name,
		PoolSize:     mainCfg.PoolSize,
		WriteTimeout: mainCfg.WriteTimeout,
		ReadTimeout:  mainCfg.ReadTimeout,
		MaxRetries:   mainCfg.MaxRetries,
	})

	var n int
	_, err := conn.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
