package postgres

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
)

// DBModel interface
type DBModel interface {
	Model(model ...interface{}) *orm.Query
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
}

// DBQuery wrap DBModel
type DBQuery struct {
	DBModel
	completed bool
}

// Rollback rollbacks query if it was transaction or returns error
func (q DBQuery) Rollback() error {
	switch t := q.DBModel.(type) {
	case *pg.Tx:
		if !q.completed {
			return t.Rollback()
		}
		return nil
	}

	return errors.New("rollback failed: not in Tx")
}

// Commit makes commit on transaction or do nothing
func (q *DBQuery) Commit() error {
	switch t := q.DBModel.(type) {
	case *pg.Tx:
		if !q.completed {
			q.completed = true

			return t.Commit()
		}
	}

	return nil
}

// NewTX returns DBQuery instance with new transaction
// DBQuery.Commit() must be called to run transaction
func (c Connector) NewTXContext(ctx context.Context) (DBQuery, error) {
	tx, err := c.conn.WithContext(ctx).Begin()
	return DBQuery{DBModel: tx}, err
}

// WithTX creates new transaction, pass it to the handler.
// Handling transaction with Rollback or Commit statements depend on result of handler
func (c Connector) WithTX(fn func(DBQuery) error) error {
	return c.WithTXContext(context.Background(), fn)
}

// WithTXContext creates new transaction, pass it to the handler.
// Handling transaction with Rollback or Commit statements depend on result of handler.
// Adds context to transaction
// err should be init in return for proper shadowing
func (c Connector) WithTXContext(ctx context.Context, fn func(DBQuery) error) (err error) {
	var tx *pg.Tx

	tx, err = c.conn.WithContext(ctx).Begin()
	if err != nil {
		return errors.Wrap(err, "new tx")
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = errors.Wrap(err, errors.Wrap(errRollback, "tx rollback").Error())
			}
			return
		}

		err = tx.Commit()
		if err != nil {
			err = errors.Wrap(err, "tx commit")
		}
	}()

	return fn(DBQuery{DBModel: tx})
}

// Query returns DBQuery instance of current db pool with context
func (c Connector) QueryContext(ctx context.Context) DBQuery {
	return DBQuery{DBModel: c.conn.WithContext(ctx)}
}
