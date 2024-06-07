package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	DTOUser struct {
		ID       int64
		Login    string
		Password string
	}

	UserStorage struct {
		db *sqlx.DB
		l  *zap.Logger
	}
)

func NewUserStorage(db *sqlx.DB, l *zap.Logger) *UserStorage {
	return &UserStorage{db, l}
}

func (s *UserStorage) Add(ctx context.Context, login, pwdHash string) (int64, error) {
	// queryText := `INSERT INTO "user"(login, password) VALUES ($1, $2);`
	queryText := `
		INSERT INTO "user"(login, password)
		VALUES ($1, $2)
		RETURNING  "user".id;`

	var uid int64
	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &uid, queryText, login, pwdHash); err != nil {
			return err
		}

		// _, err := tx.ExecContext(ctx, queryText, login, pwdHash)
		// if err != nil {
		// 	return err
		// }

		// if _, err := result.RowsAffected(); err != nil {
		// 	return err
		// }

		return nil
	})

	if err != nil {
		return -1, errors.Wrap(err, "userstorage.add")
	}

	return uid, nil
}

func (s *UserStorage) Get(ctx context.Context, login string) (*DTOUser, error) {
	queryText := `SELECT id, login, password FROM "user" WHERE login = $1;`

	user := DTOUser{}

	err := WithTx(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &user, queryText, login); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &user, errors.Wrap(err, "userstorage.get")
	}

	return &user, nil
}
