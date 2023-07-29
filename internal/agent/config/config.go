package config

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
	RateLimit         int    `env:"RATE_LIMIT"`
}

func GetConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("err parse environment variable to agent config")
	}

	return cfg
}

func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&flagCfg.ReportIntervalSec, "r", 10, "frequency of sending metrics")
	flag.IntVar(&flagCfg.PollIntervalSec, "p", 2, "metric polling frequency")
	flag.StringVar(&flagCfg.SignKey, "k", "", "includes key signature using an algorithm SHA256")
	flag.IntVar(&flagCfg.RateLimit, "l", 1, "number of simultaneously requests to the server")

	flag.Parse()

	return flagCfg
}
