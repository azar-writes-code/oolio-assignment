package config

type TelemetryConfig struct {
	METRIC_ENABLED bool   `mapstructure:"METRICS_ENABLED" default:"true"`
	LOKI_URL       string `mapstructure:"LOKI_URL" default:"http://localhost:3100"`
}
