package handlers

import (
	"net/http"
	"strings"

	"metric-alert/internal/storage"
)

const (
	metricType  = 2
	metricName  = 3
	metricValue = 4
)

type MetricAlerts struct {
	metricStorage storage.MetricStorage
}

func NewMetric(metricStorage storage.MetricStorage) MetricAlerts {
	return MetricAlerts{metricStorage: metricStorage}
}

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		req := strings.Split(r.URL.Path, "/")
		metric, err := parseURL(req)
		if err != nil {
			if err == ErrEmptyMetric {
				http.Error(w, err.Error(), http.StatusNotFound)

				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		m.metricStorage.SetMetric(metric)

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}
