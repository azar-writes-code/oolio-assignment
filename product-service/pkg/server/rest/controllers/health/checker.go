package health

import (
	"net"
	"time"
)

// HealthChecker is the Strategy interface for performing a single dependency health check.
// Adding new service checks never requires modifying the HealthController (OCP).
type HealthChecker interface {
	// Name returns the human-readable label for this dependency.
	Name() string
	// Check returns "UP" if the dependency is reachable, "DOWN" otherwise.
	Check() string
}

// TCPChecker is a generic HealthChecker that validates TCP reachability.
type TCPChecker struct {
	name    string
	address string
	timeout time.Duration
}

// NewTCPChecker constructs a TCPChecker for a given service address ("host:port").
func NewTCPChecker(name, address string) *TCPChecker {
	return &TCPChecker{
		name:    name,
		address: address,
		timeout: 2 * time.Second,
	}
}

func (t *TCPChecker) Name() string { return t.name }

func (t *TCPChecker) Check() string {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return "DOWN"
	}
	conn.Close()
	return "UP"
}

// KafkaChecker checks ALL registered Kafka brokers.
// Returns "UP" only if every broker is reachable, "DEGRADED" otherwise.
// This prevents a degraded multi-broker cluster from silently passing a single-broker probe.
type KafkaChecker struct {
	brokers []string
	timeout time.Duration
}

// NewKafkaChecker constructs a KafkaChecker that validates every broker address.
func NewKafkaChecker(brokers []string) *KafkaChecker {
	return &KafkaChecker{
		brokers: brokers,
		timeout: 2 * time.Second,
	}
}

func (k *KafkaChecker) Name() string { return "kafka" }

func (k *KafkaChecker) Check() string {
	for _, broker := range k.brokers {
		conn, err := net.DialTimeout("tcp", broker, k.timeout)
		if err != nil {
			return "DOWN"
		}
		conn.Close()
	}
	return "UP"
}

