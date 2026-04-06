package config

type KafkaConfig struct {
	BROKERS []string `mapstructure:"BROKERS"`
	GROUPID string   `mapstructure:"GROUP_ID"`
}
