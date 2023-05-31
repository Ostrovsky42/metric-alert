package main

import (
	"html/template"
	"metric-alert/internal/logger"
	"net/http"

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
	fileStorage, err := storage.NewFileRecorder(cfg.FileStoragePath, cfg.StoreIntervalSec, cfg.Restore, memStorage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error create fileRecorder")
	}

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
	go a.fileStorage.Run()

	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error start serve")
	}
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
