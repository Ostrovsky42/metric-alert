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
	metric storage.MetricStorage
}

func NewMetric(metric storage.MetricStorage) MetricAlerts {
	return MetricAlerts{metric: metric}
}

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		req := strings.Split(r.URL.Path, "/")
		if len(req) != 5 {
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		switch req[metricType] {
		case "gauge":
			name, value, err := prepareGauge(req)
			if err != nil {
				http.Error(w, "err prepare gauge", http.StatusBadRequest)

				return
			}
			m.metric.SetGauge(name, value)

		case "counter":
			name, value, err := prepareCounter(req)
			if err != nil {
				http.Error(w, "err prepare counter", http.StatusBadRequest)

				return
			}
			m.metric.Count(name, value)
		default:
			http.Error(w, "unknown metric type", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}
