package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPLogger logs each incoming request with method, path, status code, and latency.
// Follows the same slog format as the rest of the application.
func HTTPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Skip logging for successful metrics scrapes to reduce noise
		if strings.Contains(path, "metrics") && status == 200 {
			return
		}

		// Choose log level based on status code
		switch {
		case status >= 500:
			slog.Error("HTTP", "method", method, "path", path, "status", status, "latency", latency.String())
		case status >= 400:
			slog.Warn("HTTP", "method", method, "path", path, "status", status, "latency", latency.String())
		default:
			slog.Info("HTTP", "method", method, "path", path, "status", status, "latency", latency.String())
		}
	}
}
