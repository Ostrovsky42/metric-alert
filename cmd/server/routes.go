package main

import (
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"metric-alert/internal/handlers"
	"metric-alert/internal/midleware"
)

func NewRoutes(metric handlers.MetricAlerts, log zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(midleware.NewLogWriter(log).WithLogging)

	r.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, metric.UpdateMetric)
	r.Get(`/value/{metric_type}/{metric_name}`, metric.GetValue)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}
