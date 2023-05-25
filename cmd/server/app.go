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
	metric     handlers.MetricAlerts
	serverHost string
	log        zerolog.Logger
}

func NewApp(metric storage.MetricStorage, serverHost string, log zerolog.Logger) Application {
	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while parse web templates")
	}

	return Application{
		metric:     handlers.NewMetric(metric, tmp, log),
		serverHost: serverHost,
		log:        log,
	}
}

func (a Application) Run() {
	err := http.ListenAndServe(a.serverHost, NewRoutes(a.metric, a.log))
	if err != nil {
		panic(err)
	}
}
