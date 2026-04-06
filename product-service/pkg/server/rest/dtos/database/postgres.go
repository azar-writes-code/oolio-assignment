package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a production-ready pgxpool.Pool from the application DatabaseConfig.
// It pings the database after connecting to fail fast if the host is unreachable.
func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.HOST, cfg.PORT, cfg.USER, cfg.PASSWORD, cfg.NAME,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, apperrors.NewInternal("database: parse config", err)
	}

	// Pool tuning — sensible production defaults
	poolCfg.MaxConns = 25
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 1 * time.Minute
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		types := []string{"product_image", "_product_image"}
		for _, t := range types {
			dt, err := conn.LoadType(ctx, t)
			if err != nil {
				// If migrations haven't run yet, the type might not exist.
				continue
			}
			conn.TypeMap().RegisterType(dt)
		}
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, apperrors.NewInternal("database: create pool", err)
	}

	// Fail fast — don't start the server with a dead DB
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, apperrors.NewInternal("database: ping failed", err)
	}

	return pool, nil
}

func Migrate(ctx context.Context, cfg config.DatabaseConfig) error {
	pool, err := NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	dir := "pkg/server/rest/dtos/db/migrations"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return apperrors.NewInternal("read migrations directory", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return apperrors.NewInternal(fmt.Sprintf("read migration %s", entry.Name()), err)
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			return apperrors.NewInternal(fmt.Sprintf("execute migration %s", entry.Name()), err)
		}

		slog.Info("applied migration", "file", entry.Name())
	}

	return nil
}

func DeleteAllTables(ctx context.Context, cfg config.DatabaseConfig) error {
	pool, err := NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, `DO $$ DECLARE r RECORD; BEGIN FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE'; END LOOP; END $$;`)
	if err != nil {
		return apperrors.NewInternal("drop all tables", err)
	}

	return nil
}