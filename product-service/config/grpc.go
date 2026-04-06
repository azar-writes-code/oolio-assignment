package config

type GrpcConfig struct {
	PORT string `mapstructure:"PORT" default:"50051"`
	HOST string `mapstructure:"HOST" default:"0.0.0.0"`
}
