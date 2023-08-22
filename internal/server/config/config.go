// Пакет config предоставляет настройки конфигурации для сервера.
package config

import (
	"flag"

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
)

// Config представляет конфигурацию сервера.
type Config struct {
	ServerHost       string `env:"ADDRESS"`           // ServerHost Адрес работы сервера.
	StoreIntervalSec int    `env:"STORE_INTERVAL"`    // StoreIntervalSec Интервал записи метрик.
	FileStoragePath  string `env:"FILE_STORAGE_PATH"` // FileStoragePath Путь к файлу для записи метрик.
	DataBaseDSN      string `env:"DATABASE_DSN"`      // DataBaseDSN Строка подключения к базе данных.
	Restore          bool   `env:"RESTORE"`           // Restore Загружать ли сохраненные значения из файла при запуске.
	SignKey          string `env:"KEY"`               // SignKey Ключ для подписи с использованием алгоритма SHA256.
}

// GetConfig возвращает настройки сервера, считываемые из флагов командной строки и переменных окружения.
// Если переменные не предоставлены, будут использованы значения по умолчанию.
func GetConfig() Config {
	cfg := parseFlags()
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Fatal().Msg("ошибка при разборе переменных окружения для конфигурации сервера")
	}

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

	flag.Parse()

	return flagCfg
}
