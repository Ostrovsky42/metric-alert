package main

import "github.com/caarlos0/env/v6"

type Config struct {
	ServerHost        string `env:"ADDRESS"`
	ReportIntervalSec int    `env:"REPORT_INTERVAL"`
	PollIntervalSec   int    `env:"POLL_INTERVAL"`
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
	if cfg.ReportIntervalSec == 0 {
		cfg.ReportIntervalSec = reportIntervalSec
	}
	if cfg.PollIntervalSec == 0 {
		cfg.PollIntervalSec = pollIntervalSec
	}

	return cfg, nil
}
