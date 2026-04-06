package middleware

import (
	"github.com/azar-writes-code/oolio-products-backend/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(router *gin.Engine, m *telemetry.Metrics) {
	router.Use(CORSMiddleware())
	router.Use(HTTPLogger())
	router.Use(PrometheusMiddleware(m))
}
