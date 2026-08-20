package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/services/rewards/config"
)

func ConnectDB(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

func RunMigrations(dsn string, migrationsPath string) error {
	log.Println("Running rewards database migrations...")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Rewards migrations are already up to date")
			return nil
		}
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Rewards migrations applied successfully")
	return nil
}
