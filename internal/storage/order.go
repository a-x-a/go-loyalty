package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

type (
	DTOOrder struct {
		Number     string    `db:"number"`
		Status     int       `db:"status"`
		Accrual    float64   `db:"accrual"`
		UploadedAt time.Time `db:"uploaded_at"`
	}

	DTOOrders []DTOOrder

	DTOAccrualOrder struct {
		UID    int64  `db:"user_id"`
		Number string `db:"number"`
	}

	DTOAccrualOrders []DTOAccrualOrder

	OrderStorage struct {
		db *sqlx.DB
		l  *zap.Logger
	}
)

func NewOrderStorage(db *sqlx.DB, l *zap.Logger) *OrderStorage {
	return &OrderStorage{db, l}
}

func (s *OrderStorage) Add(ctx context.Context, uid int64, number string) error {
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		// проверим кому принадлежит заказ.
		var id int64
		queryText := `SELECT user_id FROM "order" WHERE number = $1;`
		if err := tx.GetContext(ctx, &id, queryText, number); err == nil {
			if id == uid {
				return customerrors.ErrOrderUploadedByUser
			}

			return customerrors.ErrOrderUploadedByAnotherUser
		}

		queryText = `INSERT INTO "order"(number, user_id) VALUES ($1, $2);`
		if _, err := tx.ExecContext(ctx, queryText, number, uid); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "orderstorage.add")
	}

	return nil
}

func (s *OrderStorage) GetAll(ctx context.Context, uid int64) (*DTOOrders, error) {
	queryText := `SELECT number, status, accrual, uploaded_at FROM "order" WHERE user_id = $1;`

	orders := DTOOrders{}
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &orders, queryText, uid); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &orders, errors.Wrap(err, "orderstorage.getall")
	}

	return &orders, nil
}

func (s *OrderStorage) GetToProcessing(ctx context.Context) (*DTOAccrualOrders, error) {
	queryText := `SELECT user_id, number
		FROM "order"
		WHERE status < 3;`
	orders := DTOAccrualOrders{}
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &orders, queryText); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &orders, errors.Wrap(err, "orderstorage.getall")
	}

	return &orders, nil
}

func (s *OrderStorage) Update(ctx context.Context, order DTOOrder) error {
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		queryText := `UPDATE "order"
			SET status = $2, accrual = $3
			WHERE number = $1;`
		if _, err := tx.ExecContext(ctx, queryText, order.Number, order.Status, order.Accrual); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "orderstorage.add")
	}

	return nil
}
