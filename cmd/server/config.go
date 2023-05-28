package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerHost       string `env:"ADDRESS"`
	StoreIntervalSec int    `env:"STORE_INTERVAL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	Restore          bool   `env:"RESTORE"`
}

func getConfig() (Config, error) {
	cfg := parseFlags()

	if serverHost, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerHost = serverHost
	}
	if storeIntervalSec, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if intInterval, err := strconv.Atoi(storeIntervalSec); err == nil {
			cfg.StoreIntervalSec = intInterval
		}
	}
	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = fileStoragePath
	}
	if restore, ok := os.LookupEnv("RESTORE"); ok {
		switch restore {
		case "true":
			cfg.Restore = true
		case "false":
			cfg.Restore = false
		}
	}

	return cfg, nil
}

func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", "localhost:8080", "server endpoint host")
	flag.IntVar(&flagCfg.StoreIntervalSec, "i", 300, "interval for writing server readings to disk")
	flag.StringVar(&flagCfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path to the file for recording readings")
	flag.BoolVar(&flagCfg.Restore, "r", true, "load saved values from the specified file at startup")

	flag.Parse()

	return flagCfg
}
