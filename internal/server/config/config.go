// Package config provides configuration settings for the application.
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"metric-alert/internal/server/logger"
)

const (
	DefaultServerHost       = "localhost:8080"
	DefaultStoreIntervalSec = 300
	DefaultFileStoragePath  = "/tmp/metrics-db.json"
	DefaultDataBaseDSN      = ""
	DefaultRestore          = true
	DefaultSignKey          = ""
)

// Config represents the application configuration.
type Config struct {
	ServerHost       string `env:"ADDRESS"`           // Address of the server endpoint.
	StoreIntervalSec int    `env:"STORE_INTERVAL"`    // Interval for writing server readings to disk.
	FileStoragePath  string `env:"FILE_STORAGE_PATH"` // Path to the file for recording readings.
	DataBaseDSN      string `env:"DATABASE_DSN"`      // Address of the database connection.
	Restore          bool   `env:"RESTORE"`           // Whether to load saved values from a file at startup.
	SignKey          string `env:"KEY"`               // Key for signature using the SHA256 algorithm.
}

// GetConfig provides configuration set through flags or environment variables.
// If no variables are provided, default values will be used.
func GetConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("error parsing environment variables to server config")
	}

	return cfg
}

// parseFlags parses command-line flags and returns a Config.
func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", DefaultServerHost, "server endpoint host")
	flag.IntVar(&flagCfg.StoreIntervalSec, "i", DefaultStoreIntervalSec, "interval for writing server readings to disk")
	flag.StringVar(&flagCfg.FileStoragePath, "f", DefaultFileStoragePath, "path to the file for recording readings")
	flag.StringVar(&flagCfg.DataBaseDSN, "d", DefaultDataBaseDSN, "string with the address of the connection to the database")
	flag.BoolVar(&flagCfg.Restore, "r", DefaultRestore, "load saved values from the specified file at startup")
	flag.StringVar(&flagCfg.SignKey, "k", DefaultSignKey, "includes key signature using an algorithm SHA256")

	flag.Parse()

	return flagCfg
}
