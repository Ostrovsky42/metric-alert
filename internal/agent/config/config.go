// Package config предоставляет настройки конфигурации для агента.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"metric-alert/internal/server/logger"

	"github.com/caarlos0/env/v6"
)

// Константы для конфигурации по умолчанию.
const (
	DefaultServerHost        = "localhost:8080"
	DefaultProfilerHost      = "localhost:6060"
	DefaultReportIntervalSec = 10
	DefaultPollIntervalSec   = 2
	DefaultSignKey           = ""
	DefaultRateLimit         = 1
	DefaultPath              = ""
	DefaultHTTP              = false
)

// Config содержит настройки агента.
type Config struct {
	LocalIP           string // LocalIP заполняется при старте программы.
	ServerHost        string `json:"server_host" env:"ADDRESS"`                 // ServerHost определяет адрес сервера.
	ProfilerHost      string `json:"profiler_host" env:"PROFILER_HOST"`         // ProfilerHost Порт на котором будет запускаться сервер для профилирования.
	ReportIntervalSec int    `json:"report_interval_sec" env:"REPORT_INTERVAL"` // ReportIntervalSec определяет интервал отправки метрик.
	PollIntervalSec   int    `json:"poll_interval_sec" env:"POLL_INTERVAL"`     // PollIntervalSec определяет интервал опроса метрик.
	SignKey           string `json:"sign_key" env:"KEY"`                        // SignKey определяет ключ подписи.
	RateLimit         int    `json:"rate_limit" env:"RATE_LIMIT"`               // RateLimit определяет ограничение скорости запросов к серверу.
	CryptoKey         string `json:"crypto_key" env:"CRYPTO_KEY"`               // CryptoKey Публичный ключ для использования асиметричного шифрования.
	JSONConfig        string `json:"json_config" env:"CONFIG"`                  // JSONConfig Путь к файлу конфигурацияй в формате JSON (самй низкий приоритет)
	IsHTTP            bool   `json:"http" env:"HTTP"`                           // IsHTTP Запуск HTTP клиента,по умолчанию запускается GRPC

}

// GetConfig возвращает настройки агента, считываемые из флагов командной строки и переменных окружения.
// Если переменные не предоставлены, будут использованы значения по умолчанию.
func GetConfig() *Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err parse environment variable to agent config")
	}

	CheckJSONConfig(&cfg)
	cfg.LocalIP, err = getLocalIP()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err get local ip")
	}

	return &cfg
}

// parseFlags разбирает флаги командной строки и возвращает настройки агента.
func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", DefaultServerHost, "server endpoint address")
	flag.StringVar(&flagCfg.ProfilerHost, "pp", DefaultProfilerHost, "profiler endpoint address")
	flag.IntVar(&flagCfg.ReportIntervalSec, "r", DefaultReportIntervalSec, "frequency of sending metrics")
	flag.IntVar(&flagCfg.PollIntervalSec, "p", DefaultPollIntervalSec, "metric polling frequency")
	flag.StringVar(&flagCfg.SignKey, "k", DefaultSignKey, "includes key signature using an algorithm SHA256")
	flag.IntVar(&flagCfg.RateLimit, "l", DefaultRateLimit, "number of simultaneously requests to the server")
	flag.StringVar(&flagCfg.CryptoKey, "crypto-key", DefaultPath, "path to the file with the public key for asymmetric encryption")
	flag.StringVar(&flagCfg.JSONConfig, "config", DefaultPath, "path to the configuration file in JSON format")
	flag.BoolVar(&flagCfg.IsHTTP, "http", DefaultHTTP, "starting the HTTP client, GRPC is started by default")

	flag.Parse()

	return flagCfg
}

// CheckJSONConfig проверяет путь к файлу хранящему конфигурацию в формате JSON
// и если он был передан через флаг или переменную окружения, заполнит не переданные значения данными из файла.
func CheckJSONConfig(cfg *Config) {
	if len(cfg.JSONConfig) != 0 {
		jsonCfg, err := readJSONConfig(cfg.JSONConfig)
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("error read JSON config")
		}
		setJSONConfig(cfg, jsonCfg)
	}
}

func readJSONConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func setJSONConfig(config *Config, jsonConfig Config) {
	if config.ServerHost == DefaultServerHost {
		config.ServerHost = jsonConfig.ServerHost
	}
	if config.ReportIntervalSec == DefaultReportIntervalSec {
		config.ReportIntervalSec = jsonConfig.ReportIntervalSec
	}
	if config.PollIntervalSec == DefaultPollIntervalSec {
		config.PollIntervalSec = jsonConfig.PollIntervalSec
	}
	if config.RateLimit == DefaultRateLimit {
		config.RateLimit = jsonConfig.RateLimit
	}
	if config.SignKey == DefaultSignKey {
		config.SignKey = jsonConfig.SignKey
	}
	if config.CryptoKey == DefaultPath {
		config.CryptoKey = jsonConfig.CryptoKey
	}
}

func getLocalIP() (string, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp != 0 && iFace.Flags&net.FlagLoopback == 0 {
			adders, err := iFace.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range adders {
				IPNet, ok := addr.(*net.IPNet)
				if ok && !IPNet.IP.IsLoopback() && IPNet.IP.To4() != nil {
					return IPNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("err found local IP address")
}
