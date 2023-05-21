package main

import (
	"net/http"

	"github.com/rs/zerolog"
	"metric-alert/internal/handlers"
	"metric-alert/internal/storage"
)

type Application struct {
	metric     handlers.MetricAlerts
	serverHost string
	log        zerolog.Logger
}

func NewApp(metric storage.MetricStorage, serverHost string, log zerolog.Logger) Application {
	return Application{
		metric:     handlers.NewMetric(metric, log),
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
