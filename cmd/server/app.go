package main

import (
	"html/template"
	"metric-alert/internal/logger"
	"net/http"
	"time"

	"metric-alert/internal/handlers"
	"metric-alert/internal/storage"
)

const templatePath = "internal/html/templates/info_page.html"

type Application struct {
	metric      handlers.MetricAlerts
	postgres    *storage.Postgres
	fileStorage *storage.FileRecorder
	serverHost  string
}

func NewApp(cfg Config) Application {
	memStorage := storage.NewMemStore()
	fileStorage, err := storage.NewFileRecorder(cfg.FileStoragePath, memStorage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error create fileRecorder")
	}
	pg, err := storage.NewPostgresDB(cfg.DataBaseDSN)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error connect to db")
	}

	if cfg.Restore {
		fileStorage.RestoreMetrics()
	}

	go StartRecording(fileStorage, cfg.StoreIntervalSec)

	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error while parse web templates")
	}

	return Application{
		metric:      handlers.NewMetric(memStorage, pg, tmp),
		postgres:    pg,
		fileStorage: fileStorage,
		serverHost:  cfg.ServerHost,
	}
}

func (a Application) Run() {
	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error start serve")
	}
}

func (a Application) Close() {
	a.postgres.Close()
}

func StartRecording(fileStorage *storage.FileRecorder, updateInterval int) {
	interval := time.Duration(updateInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		time.Sleep(interval)
		fileStorage.RecordMetrics()
	}
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
