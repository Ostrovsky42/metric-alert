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
	fileStorage *storage.FileRecorder
	serverHost  string
}

func NewApp(cfg Config) Application {
	memStorage := storage.NewMemStore()
	fileStorage, err := storage.NewFileRecorder(cfg.FileStoragePath, memStorage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error create fileRecorder")
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
		metric:      handlers.NewMetric(memStorage, tmp),
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

func StartRecording(fileStorage *storage.FileRecorder, updateInterval int) {
	interval := time.Duration(updateInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fileStorage.RecordMetrics()
		}
	}
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
