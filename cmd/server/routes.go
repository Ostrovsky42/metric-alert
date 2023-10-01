package main

import (
	"compress/gzip"

	"github.com/go-chi/chi"

	"metric-alert/internal/server/handlers"
	"metric-alert/internal/server/middleware"
)

func NewRoutes(metric handlers.MetricAlerts, signKey string, path string) *chi.Mux {
	r := chi.NewRouter()

	zipMW := middleware.NewZipMiddleware(gzip.BestSpeed)
	hashMW := middleware.NewHashMW(signKey)
	decryptorMW := middleware.NewDecryptorMW(path)

	r.Use(middleware.WithLogging, zipMW.UnZip, decryptorMW.Decrypt, hashMW.Hash, zipMW.Zip)

	r.Post(`/update/`, metric.UpdateMetricWithBody)
	r.Post(`/updates/`, metric.UpdateMetricsWithBody)
	r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)

	r.Post(`/value/`, metric.GetValueWithBody)
	r.Get(`/value/{type}/{name}`, metric.GetValue)

	r.Get("/ping", metric.PingDB)
	r.Get(`/`, metric.InfoPage)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}
