package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/database"
)

func Migrate() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// 1. Gather Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err.Error())
		os.Exit(1)
	}

	// 2. Run Database Migrations
	if err := database.Migrate(ctx, cfg.Database); err != nil {
		slog.Error("Failed to run database migrations", "error", err.Error())
		os.Exit(1)
	}
	
	slog.Info("Database migrations applied successfully")
}

func Drop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Gather Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err.Error())
		os.Exit(1)
	}

	// 2. Run Database Migrations
	if err := database.DeleteAllTables(ctx, cfg.Database); err != nil {
		slog.Error("Failed to drop database tables", "error", err.Error())
		os.Exit(1)
	}
	
	slog.Info("Database tables dropped successfully")
}