package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool wraps a pool of connections and manages communication with the database.
type DBPool struct {
	pool *pgxpool.Pool
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get new migrate instance: %w", err)
	}
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}

	return nil
}

func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	const (
		minConns = 1
		maxConns = 5
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	poolCfg.MinConns = minConns
	poolCfg.MaxConns = maxConns
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init connection pool: %w", err)
	}

	return pool, nil
}

// NewPool runs migrations and initialized new connection pool.
func NewPool(ctx context.Context, dsn string) (*DBPool, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := initPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init pool: %w", err)
	}

	return &DBPool{pool: pool}, nil
}
