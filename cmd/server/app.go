package main

import (
	"metric-alert/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"metric-alert/internal/storage"
)

type Application struct {
	metric handlers.MetricAlerts
}

func NewApp(metric storage.MetricStorage) Application {
	return Application{metric: handlers.NewMetric(metric)}
}

func (a Application) Run() {
	r := chi.NewRouter()
	r.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, a.metric.UpdateMetric)
	r.Get(`/value/{metric_type}/{metric_name}`, a.metric.GetValue)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
