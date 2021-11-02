package postgres

import (
	"github.com/erp/api/src/config"
	"github.com/go-pg/pg"
)

// Connector pg connection
type Connector struct {
	conn *pg.DB
}

// New returns postgres Connector
func New() *Connector {
	return &Connector{}
}

// NewConn init, overwrite connection
func (c *Connector) NewConn(mainCfg *config.Postgres) (err error) {
	c.conn, err = newConn(mainCfg)
	return
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

// SetConn overwrites connection
func (c *Connector) SetConn(conn *pg.DB) {
	c.conn = conn
}
