package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/handlers/validator"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/storage"
)

func (m MetricAlerts) GetValueWithBody(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = validator.ValidateGetWithBody(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metric")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	receivedMetric, err := m.metricStorage.GetMetric(r.Context(), metric.ID)
	if err != nil {
		logger.Log.Warn().Interface("metric", metric).Msg("error get metric")
		if err.Error() == storage.NotFound {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(receivedMetric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err encode data")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (m MetricAlerts) GetValue(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}

	metric.MType = chi.URLParam(r, metricType)
	metric.ID = chi.URLParam(r, metricName)

	err := validator.ValidateGet(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metric")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	receivedMetric, err := m.metricStorage.GetMetric(r.Context(), metric.ID)
	if err != nil {
		logger.Log.Warn().Interface("metric", metric).Msg("error get metric")
		if err.Error() == storage.NotFound {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	sendOK(w, *receivedMetric)
}

func (m MetricAlerts) InfoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics, err := m.metricStorage.GetAllMetric(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("error get metrics")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err = m.tmp.Execute(w, metrics); err != nil {
		logger.Log.Error().Err(err).Msg("err Execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func sendOK(w http.ResponseWriter, metric entities.Metrics) {
	w.WriteHeader(http.StatusOK)
	if metric.MType == entities.Gauge {
		fmt.Fprintf(w, "%v", *metric.Value)

		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", *metric.Delta)
}

func (m MetricAlerts) PingDB(w http.ResponseWriter, r *http.Request) {
	if err := m.metricStorage.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
