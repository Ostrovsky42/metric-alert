package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost        string `env:"ADDRESS"`
	ReportIntervalSec int    `env:"REPORT_INTERVAL"`
	PollIntervalSec   int    `env:"POLL_INTERVAL"`
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
	if cfg.ReportIntervalSec == 0 {
		cfg.ReportIntervalSec = flagCfg.ReportIntervalSec
	}
	if cfg.PollIntervalSec == 0 {
		cfg.PollIntervalSec = flagCfg.PollIntervalSec
	}

	return cfg, nil
}

func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&flagCfg.ReportIntervalSec, "r", 10, "frequency of sending metrics")
	flag.IntVar(&flagCfg.PollIntervalSec, "p", 2, "metric polling frequency")

	flag.Parse()

	return flagCfg
}
