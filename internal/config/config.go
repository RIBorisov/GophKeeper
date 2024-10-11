package config

import (
	"github.com/caarlos0/env/v11"

	"github.com/RIBorisov/GophKeeper/internal/log"
)

type AppConfig struct {
	Addr  string `env:"SERVER_ADDRESS" envDefault:":50051"`
	PgDSN string `env:"POSTGRES_DSN" envDefault:"postgresql://admin:password@localhost:5432/gophkeeper?sslmode=disable"`
}

type ServiceConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:""`
}

type Config struct {
	App     AppConfig
	Service ServiceConfig
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Failed to parse config")
	}

	return cfg
}
