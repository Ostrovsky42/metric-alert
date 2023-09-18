// Package config предоставляет настройки конфигурации для сервера.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"

	"metric-alert/internal/server/logger"
)

// Константы для конфигурации по умолчанию.
const (
	DefaultServerHost       = "localhost:8080"
	DefaultStoreIntervalSec = 300
	DefaultFileStoragePath  = "/tmp/metrics-db.json"
	DefaultDataBaseDSN      = ""
	DefaultRestore          = true
	DefaultSignKey          = ""
	DefaultPath             = ""
)

// Config представляет конфигурацию сервера.
type Config struct {
	ServerHost       string `json:"server_host" env:"ADDRESS"`                 // ServerHost Адрес работы сервера.
	StoreIntervalSec int    `json:"store_interval_sec" env:"STORE_INTERVAL"`   // StoreIntervalSec Интервал записи метрик.
	FileStoragePath  string `json:"file_storage_path" env:"FILE_STORAGE_PATH"` // FileStoragePath Путь к файлу для записи метрик.
	DataBaseDSN      string `json:"data_base_dsn" env:"DATABASE_DSN"`          // DataBaseDSN Строка подключения к базе данных.
	Restore          bool   `json:"restore" env:"RESTORE"`                     // Restore Загружать ли сохраненные значения из файла при запуске.
	SignKey          string `json:"sign_key" env:"KEY"`                        // SignKey Ключ для подписи с использованием алгоритма SHA256.
	CryptoKey        string `json:"crypto_key" env:"CRYPTO_KEY"`               // CryptoKey Приватный ключ для использования асиметричного шифрования.
	JSONConfig       string `json:"json_config" env:"CONFIG"`                  // JSONConfig Путь к файлу конфигурацияй в формате JSON (самй низкий приоритет)
}

// GetConfig возвращает настройки сервера, считываемые из флагов командной строки и переменных окружения.
// Если переменные не предоставлены, будут использованы значения по умолчанию.
func GetConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("ошибка при разборе переменных окружения для конфигурации сервера")
	}
	CheckJSONConfig(&cfg)

	return cfg
}

// parseFlags разбирает флаги командной строки и возвращает  настройки сервера.
func parseFlags() Config {
	flagCfg := Config{}
	flag.StringVar(&flagCfg.ServerHost, "a", DefaultServerHost, "хост конечной точки сервера")
	flag.IntVar(&flagCfg.StoreIntervalSec, "i", DefaultStoreIntervalSec, "интервал записи показаний сервера на диск")
	flag.StringVar(&flagCfg.FileStoragePath, "f", DefaultFileStoragePath, "путь к файлу для записи показаний")
	flag.StringVar(&flagCfg.DataBaseDSN, "d", DefaultDataBaseDSN, "строка с адресом подключения к базе данных")
	flag.BoolVar(&flagCfg.Restore, "r", DefaultRestore, "загружать сохраненные значения из указанного файла при запуске")
	flag.StringVar(&flagCfg.SignKey, "k", DefaultSignKey, "ключ для подписи с использованием алгоритма SHA256")
	flag.StringVar(&flagCfg.CryptoKey, "crypto-key", DefaultPath, "путь к файлу с приватным ключом для ассимитричного шифрования")
	flag.StringVar(&flagCfg.JSONConfig, "config", DefaultPath, "путь к файлу конфигурацияй в формате JSON")

	flag.Parse()

	return flagCfg
}

// CheckJSONConfig проверяет путь к файлу хранящему конфигурацию в формате JSON
// и если он был передан через флаг или переменную окружения, заполнит не переданные значения данными из файла.
func CheckJSONConfig(cfg *Config) {
	if len(cfg.JSONConfig) == 0 {
		jsonCfg, err := readJSONConfig(cfg.JSONConfig)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error read JSON config")
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
	if config.StoreIntervalSec == DefaultStoreIntervalSec {
		config.StoreIntervalSec = jsonConfig.StoreIntervalSec
	}
	if config.FileStoragePath == DefaultFileStoragePath {
		config.FileStoragePath = jsonConfig.FileStoragePath
	}
	if config.DataBaseDSN == DefaultDataBaseDSN {
		config.DataBaseDSN = jsonConfig.DataBaseDSN
	}
	if config.SignKey == DefaultSignKey {
		config.SignKey = jsonConfig.SignKey
	}
	if config.CryptoKey == DefaultPath {
		config.CryptoKey = jsonConfig.CryptoKey
	}
	if config.Restore {
		config.Restore = jsonConfig.Restore
	}
}
