package main

import (
	"net/http"

	"metric-alert/internal/handlers"
	"metric-alert/internal/routes"
	"metric-alert/internal/storage"
)

type Application struct {
	metric     handlers.MetricAlerts
	serverHost string
}

func NewApp(metric storage.MetricStorage, serverHost string) Application {
	return Application{
		metric:     handlers.NewMetric(metric),
		serverHost: serverHost,
	}
}

func (a Application) Run() {
	err := http.ListenAndServe(a.serverHost, routes.NewRoutes(a.metric))
	if err != nil {
		panic(err)
	}
}
