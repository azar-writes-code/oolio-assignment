package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server defines an abstraction for backend server protocols like REST or gRPC.
// By depending on this interface, the App struct achieves Dependency Inversion (DIP).
// Adding new servers never requires modifying App logic (OCP).
type Server interface {
	Start() error
	// Shutdown gracefully drains in-flight connections within the provided context deadline.
	Shutdown(ctx context.Context) error
}

// App encapsulates the entire application orchestration (SRP).
type App struct {
	servers         []Server
	shutdownTimeout time.Duration
}

// NewApp creates a new App, injecting server dependencies.
func NewApp(servers ...Server) *App {
	return &App{
		servers:         servers,
		shutdownTimeout: 10 * time.Second,
	}
}

// Run starts all registered servers and blocks until a SIGINT/SIGTERM is received
// or a server error occurs. On shutdown signal, it drains all servers gracefully.
func (a *App) Run() error {
	errChan := make(chan error, len(a.servers))

	for _, s := range a.servers {
		go func(srv Server) {
			errChan <- srv.Start()
		}(s)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		slog.Info("Shutdown signal received, draining connections...", "timeout", a.shutdownTimeout.String())
		ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
		defer cancel()

		for _, s := range a.servers {
			if err := s.Shutdown(ctx); err != nil {
				slog.Error("Server did not shut down cleanly", "error", err.Error())
			}
		}

		slog.Info("Graceful shutdown complete.")
		return nil

	case err := <-errChan:
		return err
	}
}