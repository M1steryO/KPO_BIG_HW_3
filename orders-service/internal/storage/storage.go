package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"orders-service/internal/config"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrOrderNotFound = errors.New("order not found")
)

func NewPostgresOrderStorage(dbCfg config.DBConfig) (OrderStorage, error) {
	db, err := sql.Open("postgres", dbCfg.URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	db.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)

	m, err := migrate.New(
		"file://app/migrations",
		dbCfg.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}
	fmt.Printf("storage initialized\n")
	return &Storage{db: db}, nil
}
