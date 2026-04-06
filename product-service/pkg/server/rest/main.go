package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/middleware"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/routes"
	"github.com/azar-writes-code/oolio-products-backend/pkg/telemetry"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server implements the app.Server interface for the REST protocol (SRP).
type Server struct {
	port       string
	cfg        *config.Config
	router     *gin.Engine
	metrics    *telemetry.Metrics
	httpServer *http.Server
	pool       *pgxpool.Pool
	badgerDB   *badger.DB
}

// NewServer constructs a new REST Server with all injected dependencies.
func NewServer(cfg *config.Config, m *telemetry.Metrics, pool *pgxpool.Pool, badgerDB *badger.DB) *Server {
	// Set Gin to release mode explicitly to avoid global side effects in init()
	gin.SetMode(gin.ReleaseMode)

	return &Server{
		port:     cfg.Rest.PORT,
		router:   gin.New(),
		cfg:      cfg,
		metrics:  m,
		pool:     pool,
		badgerDB: badgerDB,
	}
}



// Start registers middleware + routes, then begins listening.
func (s *Server) Start() error {
	middleware.RegisterMiddlewares(s.router, s.metrics)
	routes.RegisterRoutes(s.router, s.cfg, s.pool, s.badgerDB)

	// Expose Prometheus /metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	s.httpServer = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	slog.Info("Server started", "port", s.port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown drains in-flight connections gracefully within the context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
