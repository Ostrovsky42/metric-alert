package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"metric-alert/internal/logger"
	"net/http"

	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
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

	metric, ok := m.metricStorage.GetMetric(metric.ID)
	if !ok {
		logger.Log.Warn().Interface("metric", metric).Msg("not found metric")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	data, err := json.Marshal(metric)
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

	metric, ok := m.metricStorage.GetMetric(metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(metric.ByteValue())
}

func (m MetricAlerts) InfoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if err := m.tmp.Execute(w, m.metricStorage.GetAllMetric()); err != nil {
		logger.Log.Error().Err(err).Msg("err Execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
