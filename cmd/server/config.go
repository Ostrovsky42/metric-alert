package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost string `env:"ADDRESS"`
}

func getConfig() (Config, error) {
	parseFlags()
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.ServerHost == "" {
		cfg.ServerHost = host
	}

	return cfg, nil
}
