package main

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"metric-alert/internal/handlers"
	"metric-alert/internal/midleware"
)

func NewRoutes(metric handlers.MetricAlerts, log zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	logMW := midleware.NewLogWriter(log)
	zipMW := midleware.NewZipMiddleware(log, gzip.BestSpeed)

	r.Use(logMW.WithLogging, zipMW.UnZip, zipMW.Zip)

	r.Post(`/update/`, metric.UpdateMetricWithBody)
	r.Post(`/value/`, metric.GetValueWithBody)

	r.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, metric.UpdateMetric)
	r.Get(`/value/{metric_type}/{metric_name}`, metric.GetValue)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}
