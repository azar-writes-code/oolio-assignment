package middleware

import (
	"strconv"
	"time"

	"github.com/azar-writes-code/oolio-products-backend/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware records standard HTTP metrics per request.
// Injected with *telemetry.Metrics to avoid global state (DIP).
func PrometheusMiddleware(m *telemetry.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.HTTPRequestsInFlight.Inc()
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath() // group parametrized paths like /api/v1/invoices/:id
		if path == "" {
			path = "unknown"
		}

		m.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
		m.HTTPRequestsInFlight.Dec()
	}
}
