package config

import (
	"github.com/caarlos0/env/v11"

	"github.com/RIBorisov/GophKeeper/internal/log"
)

type AppConfig struct {
	Addr        string `env:"SERVER_ADDRESS" envDefault:":50051"`
	PgDSN       string `env:"POSTGRES_DSN" envDefault:"postgresql://admin:password@localhost:5432/gophkeeper"`
	CertPath    string `env:"CERT_PATH" envDefault:"tls/server.crt"`
	CertKeyPath string `env:"CERT_KEY_PATH" envDefault:"tls/server.key"`
}

type ServiceConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"1Kg7nVcKla09d2Hf"`
}

type S3Config struct {
	BucketName      string `env:"S3_BUCKET_NAME" envDefault:"bucket"`
	Endpoint        string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	AccessKeyID     string `env:"S3_AK_ID" envDefault:"admin"`
	SecretAccessKey string `env:"S3_SECRET_AK" envDefault:"password"`
}

type Config struct {
	App     AppConfig
	Service ServiceConfig
	S3      S3Config `envPrefix:"S3"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Failed to parse config")
	}

	return cfg
}
