package handlers

import (
	"github.com/go-chi/chi"
	"html/template"
	"net/http"

	"github.com/rs/zerolog"
	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/helpers"
	"metric-alert/internal/storage"
)

type MetricAlerts struct {
	metricStorage storage.MetricStorage
	tmp           *template.Template
	log           zerolog.Logger
}

func NewMetric(metricStorage storage.MetricStorage, tmp *template.Template, log zerolog.Logger) MetricAlerts {
	return MetricAlerts{
		metricStorage: metricStorage,
		tmp:           tmp,
		log:           log,
	}
}

func (m MetricAlerts) UpdateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := helpers.UnmarshalBody(r.Body, &metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	err = validator.ValidateUpdateWithBody(metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err validate metric")
		if err.Error() == "empty metric name" {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric = m.metricStorage.SetMetric(metric)

	data, err := helpers.EncodeData(metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err encode data")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data.Bytes())
}

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}

	metric.MType = chi.URLParam(r, "metric_type")
	metric.ID = chi.URLParam(r, "metric_name")
	mValue := chi.URLParam(r, "metric_value")

	err := validator.ValidateUpdate(&metric, mValue)
	if err != nil {
		m.log.Error().Err(err).Msg("err validate metric")
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
