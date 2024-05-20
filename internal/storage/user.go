package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// type DTOUser struct {
// 	ID       int64
// 	Login    string
// 	Password string
// }

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db}
}

func (s *UserStorage) AddUser(ctx context.Context, login, pwdHash string) error {
	queryText := `INSERT INTO user(login, password) VALUES ($1, $2);`

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		result, err := tx.ExecContext(ctx, queryText, login, pwdHash)
		if err != nil {
			return err
		}

		if _, err := result.RowsAffected(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "UserStorage.AddUser")
	}

	return nil
}

func (s *UserStorage) GetUserID(ctx context.Context, login, pwdHash string) (int64, error) {
	var uid int64

	queryText := `SELECT id FROM user WHERE login = $1 AND password = $2;`

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, uid, queryText, login, pwdHash); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return uid, errors.Wrap(err, "UserStorage.GetUserID")
	}

	return uid, nil
}
