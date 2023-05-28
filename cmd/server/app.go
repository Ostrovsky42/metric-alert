package main

import (
	"html/template"
	"net/http"

	"github.com/rs/zerolog"
	"metric-alert/internal/handlers"
	"metric-alert/internal/storage"
)

const templatePath = "internal/html/templates/info_page.html"

type Application struct {
	metric      handlers.MetricAlerts
	fileStorage *storage.FileRecorder
	serverHost  string
	log         zerolog.Logger
}

func NewApp(cfg Config, log zerolog.Logger) Application {
	memStorage := storage.NewMemStore()
	fileStorage, err := storage.NewFileRecorder(cfg.FileStoragePath, cfg.StoreIntervalSec, cfg.Restore, memStorage, log)

	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while parse web templates")
	}

	return Application{
		metric:      handlers.NewMetric(memStorage, tmp, log),
		fileStorage: fileStorage,
		serverHost:  cfg.ServerHost,
		log:         log,
	}
}

func (a Application) Run() {
	go a.fileStorage.Run()

	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric, a.log))
	if err != nil {
		a.log.Fatal().Err(err).Msg("Error start serve")
	}
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
