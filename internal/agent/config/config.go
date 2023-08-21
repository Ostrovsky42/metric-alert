package config

import (
	"flag"

	"metric-alert/internal/server/logger"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultServerHost        = "localhost:8080"
	DefaultReportIntervalSec = 10
	DefaultPollIntervalSec   = 2
	DefaultSignKey           = ""
	DefaultRateLimit         = 1
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
	flag.StringVar(&flagCfg.ServerHost, "a", DefaultServerHost, "server endpoint address")
	flag.IntVar(&flagCfg.ReportIntervalSec, "r", DefaultReportIntervalSec, "frequency of sending metrics")
	flag.IntVar(&flagCfg.PollIntervalSec, "p", DefaultPollIntervalSec, "metric polling frequency")
	flag.StringVar(&flagCfg.SignKey, "k", DefaultSignKey, "includes key signature using an algorithm SHA256")
	flag.IntVar(&flagCfg.RateLimit, "l", DefaultRateLimit, "number of simultaneously requests to the server")

	flag.Parse()

	return flagCfg
}
