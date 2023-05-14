package routes

import (
	"github.com/go-chi/chi"
	"metric-alert/internal/handlers"
)

func NewRoutes(metric handlers.MetricAlerts) *chi.Mux {
	r := chi.NewRouter()
	r.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, metric.UpdateMetric)
	r.Get(`/value/{metric_type}/{metric_name}`, metric.GetValue)

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}
