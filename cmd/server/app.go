package main

import (
	"html/template"
	"metric-alert/internal/handlers"
	"metric-alert/internal/logger"
	"metric-alert/internal/repository"
	"net/http"
)

const templatePath = "internal/html/templates/info_page.html"

type Application struct {
	metric     handlers.MetricAlerts
	storage    *repository.Repository
	serverHost string
}

func NewApp(cfg Config) Application {
	memStorage, err := repository.InitRepo(
		cfg.FileStoragePath,
		cfg.DataBaseDSN,
		cfg.StoreIntervalSec,
		cfg.Restore,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed init storage")
	}

	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error while parse web templates")
	}

	return Application{
		metric:     handlers.NewMetric(memStorage, tmp),
		storage:    memStorage,
		serverHost: cfg.ServerHost,
	}
}

func (a Application) Run() {
	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error start serve")
	}
}

func (a Application) Close() {
	a.storage.Close()
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
