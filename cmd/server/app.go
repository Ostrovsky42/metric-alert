package main

import (
	"html/template"
	"metric-alert/internal/server/config"
	"metric-alert/internal/server/handlers"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/repository"
	"net/http"
)

const templatePath = "internal/server/html/templates/info_page.html"

type Application struct {
	metric     handlers.MetricAlerts
	storage    repository.MetricRepo
	serverHost string
	signKey    string
}

func NewApp(cfg config.Config) Application {
	memRepo, err := repository.InitRepo(
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
		metric:     handlers.NewMetric(memRepo, tmp),
		storage:    memRepo,
		serverHost: cfg.ServerHost,
		signKey:    cfg.SignKey,
	}
}

func (a Application) Run() {
	a.signKey = ""
	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric, a.signKey))
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
