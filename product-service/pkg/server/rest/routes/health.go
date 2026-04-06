package routes

import (
	"fmt"

	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/controllers/health"
	"github.com/gin-gonic/gin"
)

// HealthRoutes registers liveness and readiness probes.
// Checkers are composed here — the composition root — so the controller stays closed for modification.
func HealthRoutes(rg *gin.RouterGroup, port string, cfg *config.Config) {
	checkers := buildCheckers(cfg)
	hc := health.NewHealthController(port, cfg.App.Version, cfg.App.ServiceName, checkers...)

	hrg := rg.Group("/health")
	hrg.GET("/", hc.Health)
	hrg.GET("/ping", hc.Ping)
}

// buildCheckers constructs all health checkers from configuration.
// To add a new dependency health check, add a new checker here — no other file changes needed (OCP).
func buildCheckers(cfg *config.Config) []health.HealthChecker {
	checkers := []health.HealthChecker{
		health.NewTCPChecker("database", fmt.Sprintf("%s:%s", cfg.Database.HOST, cfg.Database.PORT)),
		health.NewTCPChecker("redis", fmt.Sprintf("%s:%s", cfg.Redis.HOST, cfg.Redis.PORT)),
	}
	if len(cfg.Kafka.BROKERS) > 0 {
		checkers = append(checkers, health.NewKafkaChecker(cfg.Kafka.BROKERS))
	}
	return checkers
}
