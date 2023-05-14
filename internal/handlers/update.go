package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/storage"
	"metric-alert/internal/types"
)

type MetricAlerts struct {
	metricStorage storage.MetricStorage
}

func NewMetric(metricStorage storage.MetricStorage) MetricAlerts {
	return MetricAlerts{metricStorage: metricStorage}
}

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := types.Metric{}

	metric.MetricType = chi.URLParam(r, "metric_type")
	metric.MetricName = chi.URLParam(r, "metric_name")
	mValue := chi.URLParam(r, "metric_value")

	err := validator.ValidateUpdate(&metric, mValue)
	if err != nil {
		if err.Error() == "empty metric name" {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	m.metricStorage.SetMetric(metric)

	w.WriteHeader(http.StatusOK)
}
