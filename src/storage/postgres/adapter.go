package postgres

import (
	"context"

	log "github.com/sirupsen/logrus"
	pgv2 "gitlab.yalantis.com/gophers/pg/v2"
)

type LoggerAdapter struct {
	*log.Logger
}

func (l *LoggerAdapter) WithCtx(ctx context.Context) pgv2.Debugger {
	return l.Logger.WithContext(ctx)
}
