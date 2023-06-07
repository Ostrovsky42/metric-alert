package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"html/template"
	"metric-alert/internal/logger"
	"net/http"

	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/storage"
)

const (
	metricType  = "type"
	metricName  = "name"
	metricValue = "value"
)

type MetricAlerts struct {
	metricStorage storage.MetricStorage
	tmp           *template.Template
}

func NewMetric(metricStorage storage.MetricStorage, tmp *template.Template) MetricAlerts {
	return MetricAlerts{
		metricStorage: metricStorage,
		tmp:           tmp,
	}
}

func (m MetricAlerts) UpdateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	err = validator.ValidateUpdateWithBody(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metric")
		if err.Error() == "empty metric name" {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric, err = m.metricStorage.SetMetric(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error set metric")
		w.WriteHeader(http.StatusInternalServerError)

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

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}

	metric.MType = chi.URLParam(r, metricType)
	metric.ID = chi.URLParam(r, metricName)
	mValue := chi.URLParam(r, metricValue)

	err := validator.ValidateUpdate(&metric, mValue)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metric")
		if err.Error() == "empty metric name" {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = m.metricStorage.SetMetric(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error set metric")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
