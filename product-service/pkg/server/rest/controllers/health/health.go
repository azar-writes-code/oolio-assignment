package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health is the Readiness Probe. It iterates all registered HealthCheckers.
// New dependencies are added at the composition root — this function never changes (OCP).
func (hc *HealthController) Health(c *gin.Context) {
	dependencies := make(map[string]string, len(hc.checkers))
	overallStatus := "OK"

	for _, checker := range hc.checkers {
		status := checker.Check()
		dependencies[checker.Name()] = status
		if status == "DOWN" {
			overallStatus = "DEGRADED"
		}
	}

	result := map[string]any{
		"status": overallStatus,
		"service": map[string]any{
			"name":    hc.serviceName,
			"version": hc.version,
			"port":    hc.port,
		},
		"dependencies": dependencies,
	}

	statusCode := http.StatusOK
	if overallStatus == "DEGRADED" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, result)
}
