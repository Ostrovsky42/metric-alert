package main

import (
	"net/http"

	"metric-alert/internal/handlers"
	"metric-alert/internal/storage"
)

type Application struct {
	metric handlers.MetricAlerts
}

func NewApp(metric storage.MetricStorage) Application {
	return Application{metric: handlers.NewMetric(metric)}
}

func (a Application) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, a.metric.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
