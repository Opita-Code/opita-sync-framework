package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func contextBackground() context.Context {
	return context.Background()
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
