package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:",squash"`
	Rest      RestConfig      `mapstructure:"REST"`
	Grpc      GrpcConfig      `mapstructure:"GRPC"`
	Database  DatabaseConfig  `mapstructure:"DB"`
	Redis     RedisConfig     `mapstructure:"REDIS"`
	Kafka     KafkaConfig     `mapstructure:"KAFKA"`
	Telemetry TelemetryConfig `mapstructure:"TELEMETRY"`
	Auth      AuthConfig      `mapstructure:"AUTH"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load() // Load .env file if it exists, ignore errors if it doesn't

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Load optional config.yaml (cwd or ./config dir). Env vars always override file values.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Default values — single source of truth
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("ENVIRONMENT", "local")
	viper.SetDefault("SERVICE_NAME", "oolio-product-service")
	viper.SetDefault("VERSION", "dev")
	viper.SetDefault("REST.PORT", "8080")
	viper.SetDefault("REST.HOST", "0.0.0.0")
	viper.SetDefault("GRPC.PORT", "50051")
	viper.SetDefault("GRPC.HOST", "0.0.0.0")

	// Infrastructure defaults (match infra/docker/docker-compose.yml)
	viper.SetDefault("DB.HOST", "localhost")
	viper.SetDefault("DB.PORT", "5432")
	viper.SetDefault("DB.USER", "user")
	viper.SetDefault("DB.PASSWORD", "password")
	viper.SetDefault("DB.NAME", "oolio_product_db")
	viper.SetDefault("REDIS.HOST", "localhost")
	viper.SetDefault("REDIS.PORT", "6379")
	viper.SetDefault("REDIS.DB", 0)
	viper.SetDefault("KAFKA.BROKERS", []string{"localhost:9094"})
	viper.SetDefault("KAFKA.GROUP_ID", "oolio-product-group")
	viper.SetDefault("TELEMETRY.METRICS_ENABLED", true)
	viper.SetDefault("TELEMETRY.LOKI_URL", "http://localhost:3100")

	// Auth defaults
	viper.SetDefault("AUTH.API_KEY", "apitest")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
