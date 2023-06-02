package main

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"metric-alert/internal/handlers"
	"metric-alert/internal/midleware"
)

func NewRoutes(metric handlers.MetricAlerts) *chi.Mux {
	r := chi.NewRouter()

	zipMW := midleware.NewZipMiddleware(gzip.BestSpeed)

	r.Use(midleware.WithLogging, zipMW.UnZip, zipMW.Zip)

	r.Post(`/update/`, metric.UpdateMetricWithBody)
	r.Post(`/value/`, metric.GetValueWithBody)

	r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
	r.Get(`/value/{type}/{name}`, metric.GetValue)
	r.Get(`/`, metric.InfoPage)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}
