package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"single-window/config"
	"single-window/pkg/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
	scriptsPath     = "sql"
)

func RunMigration(cfg *config.Config, l *logger.Logger) error {
	if len(cfg.PG.URL) == 0 {
		return errors.New("migrate: environment variable not declared: PG_URL")
	}

	databaseURL := cfg.PG.URL + "?sslmode=disable"

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		l.Debug("migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return errors.New("migrate: postgres connect error: %s")
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("migrate: no change")
		return nil
	}

	if err != nil {
		return tryCloseMigrateWithError(m, err)
	}

	l.Info("Migrate: up success")

	if cfg.Mode == "local" {
		if err = executeSQLScripts(databaseURL); err != nil {
			return tryCloseMigrateWithError(m, err)
		}

		l.Info("all SQL scripts executed successfully")
	}

	return tryCloseMigrateWithError(m, nil)
}

func executeSQLScripts(databaseURL string) error {
	files, err := os.ReadDir(scriptsPath)
	if err != nil {
		return fmt.Errorf("error reading sql scripts directory: %w", err)
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			log.Printf("Executing SQL script: %s", file.Name())
			content, err := os.ReadFile(filepath.Join(scriptsPath, file.Name()))
			if err != nil {
				return tryCloseDBWithError(db, fmt.Errorf("failed to read sql file %s, err: %w", file.Name(), err))
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return tryCloseDBWithError(db, fmt.Errorf("failed to execute sql script %s, err: %w", file.Name(), err))
			}
		}
	}

	return tryCloseDBWithError(db, nil)
}

func tryCloseDBWithError(db *sql.DB, err error) error {
	if closeErr := db.Close(); closeErr != nil {
		return errors.Join(err, fmt.Errorf("failed to close db, err: %w", closeErr))
	}
	return err
}

func tryCloseMigrateWithError(m *migrate.Migrate, err error) error {
	var resultErr error
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		resultErr = fmt.Errorf("failed to close source, err: %w", sourceErr)
	}
	if databaseErr != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("failed to close database, err: %w", databaseErr))
	}
	return errors.Join(err, resultErr)
}
