package app

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func migrationRun(dsn string, log *zap.Logger) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	sourceURL := "file://migrations"
	if pwd, err := os.Getwd(); err == nil {
		sourceURL = fmt.Sprintf("file://%s/migrations", strings.ReplaceAll(pwd, "\\", "/"))
	}

	log.Info("migration", zap.String("sourceURL", sourceURL), zap.Any("driver", driver))
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "go-loyalty", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	defer m.Close()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("migrations: no change")
			return nil
		}

		return err
	}

	log.Info("migrations: success")
	return nil
}
