package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"metric-alert/internal/server/logger"
)

type Config struct {
	ServerHost        string `env:"ADDRESS"`
	ReportIntervalSec int    `env:"REPORT_INTERVAL"`
	PollIntervalSec   int    `env:"POLL_INTERVAL"`
	SignKey           string `env:"KEY"`
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
	flag.StringVar(&flagCfg.SignKey, "k", "", "includes key signature using an algorithm SHA256")

	flag.Parse()

	return flagCfg
}
