package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewConnection возвращает новое соединенние с базой данных.
func NewConnection(dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		// zap.L().Error("failed to create a database connection", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		// zap.L().Error("failed ping the database", zap.Error(err))
		return nil, err
	}

	return db, nil
}

type WithTxFn func(ctx context.Context, tx *sqlx.Tx) error

func WithTx(ctx context.Context, db *sqlx.DB, f WithTxFn) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "db.BeginTxx")
	}

	if err = f(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, "tx.Rollback")
		}

		return errors.Wrap(err, "tx.WithTxFunc")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "tx.Commit")
	}

	return nil
}
