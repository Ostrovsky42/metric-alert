package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"metric-alert/internal/types"
)

func (m MetricAlerts) GetValue(w http.ResponseWriter, r *http.Request) {
	metric := types.Metric{}

	metric.MetricType = chi.URLParam(r, "metric_type")
	metric.MetricName = chi.URLParam(r, "metric_name")

	err := ValidateGet(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	metric, ok := m.metricStorage.GetMetric(metric)
	if !ok {
		http.Error(w, "", http.StatusNotFound)

		return
	}

	if metric.MetricType == types.Gauge {
		w.Write([]byte(fmt.Sprintf("%v", metric.GaugeValue)))
	} else {
		w.Write([]byte(fmt.Sprintf("%v", metric.CounterValue)))
	}

	w.WriteHeader(http.StatusOK)
}
