package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

type (
	DTOBalance struct {
		Current   float64 `db:"current"`
		Withdrawn float64 `db:"withdrawn"`
	}

	DTOWithdrawal struct {
		Order       string    `json:"order" db:"order"`
		Sum         float64   `json:"sum,omitempty" db:"sum"`
		ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
	}

	DTOWithdrawals []DTOWithdrawal
)

type BalanceStorage struct {
	db *sqlx.DB
	l  *zap.Logger
}

func NewBalanceStorage(db *sqlx.DB, l *zap.Logger) *BalanceStorage {
	return &BalanceStorage{db, l}
}

func (s *BalanceStorage) Create(ctx context.Context, uid int64) error {
	queryText := `INSERT INTO "balance"(user_id) VALUES ($1);`

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		result, err := tx.ExecContext(ctx, queryText, uid)
		if err != nil {
			return err
		}

		if _, err := result.RowsAffected(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "balancestorage.crate")
	}

	return nil
}

func (s *BalanceStorage) Update(ctx context.Context, uid int64) error {
	return nil
}

func (s *BalanceStorage) Get(ctx context.Context, uid int64) (*DTOBalance, error) {
	queryText := `SELECT current, withdrawn FROM "balance" WHERE user_id = $1;`

	balance := DTOBalance{}
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &balance, queryText, uid); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &balance, errors.Wrap(err, "balancestorage.get")
	}

	return &balance, nil
}

func (s *BalanceStorage) Withdraw(ctx context.Context, uid int64, number string, sum float64) error {
	balance := DTOBalance{}

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		// добавим запись о списании.
		queryText := `INSERT INTO "withdraw"
			(user_id, "order", sum)
			VALUES($1, $2, $3)`

		_, err := tx.ExecContext(ctx, queryText, uid, number, sum)
		if err != nil {
			return err
		}

		// if _, err := result.RowsAffected(); err != nil {
		// 	return err
		// }

		// обновим баланс.
		queryText = `INSERT INTO "balance"
			(user_id, current, withdrawn)
			SELECT b.user_id, b.current-$2, b.withdrawn+$2
				FROM "balance" AS b
				WHERE user_id = $1
			ON CONFLICT (user_id) DO UPDATE
				SET current = balance.current - $2, withdrawn = balance.withdrawn + $2`

		_, err = tx.ExecContext(ctx, queryText, uid, sum)
		if err != nil {
			return err
		}

		// if _, err := result.RowsAffected(); err != nil {
		// 	return err
		// }

		// получим баланс.
		queryText = `SELECT current, withdrawn
			FROM "balance"
			WHERE user_id = $1;`

		if err := tx.GetContext(ctx, &balance, queryText, uid); err != nil {
			return err
		}

		// проверка баланса.
		if balance.Current < 0 {
			// не достаточно средств.
			return customerrors.ErrInsufficientFunds
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "balancestorage.withdraw")
	}

	return nil
}

func (s *BalanceStorage) GetWithdrawals(ctx context.Context, uid int64) (*DTOWithdrawals, error) {
	queryText := `SELECT "order", sum, processed_at FROM "withdraw" WHERE user_id = $1;`

	withdrawals := DTOWithdrawals{}
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &withdrawals, queryText, uid); err != nil {
			fmt.Println("err", err)
			return err
		}

		return nil
	})

	if err != nil {
		return &withdrawals, errors.Wrap(err, "balancestorage.getwithdrawals")
	}

	return &withdrawals, nil
}

func (s *BalanceStorage) Accrual(ctx context.Context, uid int64, sum float64) error {
	queryText := `INSERT INTO "balance"
		(user_id, current)
		SELECT b.user_id, b.current+$2
			FROM "balance" AS b
			WHERE user_id = $1
		ON CONFLICT (user_id) DO UPDATE
			SET current = balance.current + $2`

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		// обновим баланс.
		_, err := tx.ExecContext(ctx, queryText, uid, sum)
		if err != nil {
			return err
		}

		// if _, err := result.RowsAffected(); err != nil {
		// 	return err
		// }

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "balancestorage.withdraw")
	}

	return nil
}
