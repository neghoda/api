package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/erp/api/src/models"
)

func toServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
		return models.ErrNotFound
	}

	return err
}
