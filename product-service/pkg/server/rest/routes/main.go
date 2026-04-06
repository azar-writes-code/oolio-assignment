package routes

import (
	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes attaches all versioned route groups.
// cfg is passed rather than *RestConfig so sub-route groups can access full config (DB, Kafka etc.) for composition.
func RegisterRoutes(router *gin.Engine, cfg *config.Config, pool *pgxpool.Pool, badgerDB *badger.DB) {
	rg := router.Group("/api/v1")
	HealthRoutes(rg, cfg.Rest.PORT, cfg)
	ProductsRoutes(rg, cfg, pool)
	OrderRoutes(rg, cfg, pool, badgerDB)
}
