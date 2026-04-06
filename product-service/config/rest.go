package config

type RestConfig struct {
	PORT string `mapstructure:"PORT" default:"8080"`
	HOST string `mapstructure:"HOST" default:"0.0.0.0"`
}
