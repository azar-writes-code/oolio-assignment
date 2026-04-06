package config

type RedisConfig struct {
	HOST     string `mapstructure:"HOST"`
	PORT     string `mapstructure:"PORT"`
	PASSWORD string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}
