package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"metric-alert/internal/logger"
)

type Config struct {
	ServerHost       string `env:"ADDRESS"`
	StoreIntervalSec int    `env:"STORE_INTERVAL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN      string `env:"DATABASE_DSN"`
	Restore          bool   `env:"RESTORE"`
}

func getConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("err parse environment variable to config")
	}

	return cfg
}

func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", "localhost:8080", "server endpoint host")
	flag.IntVar(&flagCfg.StoreIntervalSec, "i", 300, "interval for writing server readings to disk")
	flag.StringVar(&flagCfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path to the file for recording readings")
	flag.StringVar(&flagCfg.DataBaseDSN, "d", "", "string with the address of the connection to the database")
	flag.BoolVar(&flagCfg.Restore, "r", true, "load saved values from the specified file at startup")

	flag.Parse()

	return flagCfg
}
