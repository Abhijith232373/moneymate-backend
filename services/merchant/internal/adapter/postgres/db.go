package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/config"
)

// ConnectDB establishes a connection pool to PostgreSQL.
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
	log.Println("Running database migrations...")

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		fallback := "./migrations"
		if _, err := os.Stat(fallback); err == nil {
			log.Printf("Migrations path %s not found, using fallback %s", migrationsPath, fallback)
			migrationsPath = fallback
		} else {
			// Try one more level up just in case it's run from cmd/
			fallback = "../migrations"
			if _, err := os.Stat(fallback); err == nil {
				log.Printf("Migrations path %s not found, using fallback %s", migrationsPath, fallback)
				migrationsPath = fallback
			}
		}
	}

	// Clean up path format for golang-migrate
	// If it's an absolute path (starts with /), file://%s will result in file:///path
	pathURL := fmt.Sprintf("file://%s", migrationsPath)
	if !strings.HasPrefix(migrationsPath, "/") {
		pathURL = fmt.Sprintf("file://%s", migrationsPath) // relative path
	}

	m, err := migrate.New(pathURL, dsn)
	if err != nil {
		return fmt.Errorf("could not create migrate instance for path %s: %w", pathURL, err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migrations are already up to date ✅")
			return nil
		}
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations applied successfully ✅")
	return nil
}
