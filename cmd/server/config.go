package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost string `env:"ADDRESS"`
}

func getConfig() (Config, error) {
	flagCfg := parseFlags()
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.ServerHost == "" {
		cfg.ServerHost = flagCfg.ServerHost
	}

	return cfg, nil
}

func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", "localhost:8080", "server endpoint host")

	flag.Parse()

	return flagCfg
}
