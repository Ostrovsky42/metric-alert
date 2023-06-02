package main

import (
	"flag"
	"metric-alert/internal/logger"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost        string `env:"ADDRESS"`
	ReportIntervalSec int    `env:"REPORT_INTERVAL"`
	PollIntervalSec   int    `env:"POLL_INTERVAL"`
}

func getConfig() (Config, error) {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("")
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
