package config

type AppConfig struct {
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	LogFormat   string `mapstructure:"LOG_FORMAT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Version     string `mapstructure:"VERSION"`
}
