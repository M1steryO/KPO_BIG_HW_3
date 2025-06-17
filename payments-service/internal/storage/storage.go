package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"payments-service/internal/config"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewPostgresAccountStorage(cfg config.DBConfig) (AccountStorage, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	m, err := migrate.New(
		"file://app/migrations",
		cfg.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}
	fmt.Printf("storage initialized\n")
	return &PostgresAccountStorage{db: db}, nil
}
