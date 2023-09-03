// Пакет config предоставляет настройки конфигурации для агента.
package config

import (
	"flag"

	"metric-alert/internal/server/logger"

	"github.com/caarlos0/env/v6"
)

// Константы для конфигурации по умолчанию.
const (
	DefaultServerHost        = "localhost:8080"
	DefaultReportIntervalSec = 10
	DefaultPollIntervalSec   = 2
	DefaultSignKey           = ""
	DefaultRateLimit         = 1
)

// Config содержит настройки агента.
type Config struct {
	ServerHost        string `env:"ADDRESS"`         // ServerHost определяет адрес сервера.
	ReportIntervalSec int    `env:"REPORT_INTERVAL"` // ReportIntervalSec определяет интервал отправки метрик.
	PollIntervalSec   int    `env:"POLL_INTERVAL"`   // PollIntervalSec определяет интервал опроса метрик.
	SignKey           string `env:"KEY"`             // SignKey определяет ключ подписи.
	RateLimit         int    `env:"RATE_LIMIT"`      // RateLimit определяет ограничение скорости запросов к серверу.
}

// GetConfig возвращает настройки агента, считываемые из флагов командной строки и переменных окружения.
// Если переменные не предоставлены, будут использованы значения по умолчанию.
func GetConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("err parse environment variable to agent config")
	}

	return cfg
}

// parseFlags разбирает флаги командной строки и возвращает настройки агента.
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
