package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/app"
	"github.com/azar-writes-code/oolio-products-backend/pkg/logger"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/database"
	"github.com/azar-writes-code/oolio-products-backend/pkg/telemetry"
)

// version is injected at build time via: go build -ldflags "-X main.version=1.2.3"
var version = "dev"

func main() {
	ctx := context.Background()

	// 1. Gather Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err.Error())
		os.Exit(1)
	}
	if cfg.App.Version == "dev" && version != "dev" {
		cfg.App.Version = version
	}

	// 2. Setup Telemetry
	var lokiWriter *telemetry.LokiWriter
	if cfg.Telemetry.METRIC_ENABLED {
		lokiWriter = telemetry.NewLokiWriter(cfg.Telemetry.LOKI_URL, cfg.App.ServiceName)
	}
	metrics := telemetry.NewMetrics(cfg.App.ServiceName)

	// 3. Setup Global Logger
	if lokiWriter != nil {
		logger.Init(cfg.App.LogLevel, cfg.App.LogFormat, lokiWriter)
	} else {
		logger.Init(cfg.App.LogLevel, cfg.App.LogFormat)
	}
	slog.Info("Starting Oolio Online Food Ordering System Backend",
		"service", cfg.App.ServiceName,
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)

	// 4. Database — connection pool + repositories
	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("Database connected", "host", cfg.Database.HOST, "name", cfg.Database.NAME)

	// 4b. BadgerDB — shared instance for coupons
	badgerDB, err := database.OpenBadger()
	if err != nil {
		slog.Error("failed to open badger database", "error", err.Error())
		os.Exit(1)
	}
	defer badgerDB.Close()
	slog.Info("BadgerDB connected")

	// 5. Instantiate Servers (Dependency Injection root)
	restServer := rest.NewServer(cfg, metrics, pool, badgerDB)

	// 6. Run App Orchestrator
	application := app.NewApp(restServer)
	if err := application.Run(); err != nil {
		slog.Error("application runtime error", "error", err.Error())
	}
}
