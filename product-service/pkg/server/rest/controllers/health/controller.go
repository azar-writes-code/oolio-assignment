package health

// HealthController handles liveness and readiness probes.
// It depends on abstractions ([]HealthChecker) rather than concrete config types (DIP).
type HealthController struct {
	port        string
	version     string
	serviceName string
	checkers    []HealthChecker
}

// NewHealthController constructs a HealthController with injected dependencies.
// Callers register HealthCheckers at the composition root — the controller never needs changing (OCP).
func NewHealthController(port, version, serviceName string, checkers ...HealthChecker) *HealthController {
	return &HealthController{
		port:        port,
		version:     version,
		serviceName: serviceName,
		checkers:    checkers,
	}
}
