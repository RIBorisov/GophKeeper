package config

type AppConfig struct {
}

type ServiceConfig struct {
}

type Config struct {
	App     AppConfig
	Service ServiceConfig
}

func Load() *Config {
	var cfg *Config
	
	return cfg
}
