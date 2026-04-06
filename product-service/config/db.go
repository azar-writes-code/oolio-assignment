package config

type DatabaseConfig struct {
	HOST     string `mapstructure:"HOST"`
	PORT     string `mapstructure:"PORT"`
	USER     string `mapstructure:"USER"`
	PASSWORD string `mapstructure:"PASSWORD"`
	NAME     string `mapstructure:"NAME"`
}
