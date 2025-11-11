package config

import (
	"os"
)

type InternalWsConfig struct {
	srvHost   string
	srvPort   string
	jwtSecret string
}

func NewInternalWsConfig() *InternalWsConfig {
	return &InternalWsConfig{
		srvHost:   "localhost",
		srvPort:   "7979",
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

func (cfg *InternalWsConfig) GetServerHost() string {
	return cfg.srvHost
}

func (cfg *InternalWsConfig) GetServerPort() string {
	return cfg.srvPort
}

func (cfg *InternalWsConfig) GetJwtSecret() string {
	return cfg.jwtSecret
}
