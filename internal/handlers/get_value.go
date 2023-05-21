package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"

	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/helpers"
)

func (m MetricAlerts) GetValueWithBody(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := helpers.UnmarshalBody(r.Body, &metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = validator.ValidateGetWithBody(metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err validate metric")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric, ok := m.metricStorage.GetMetric(metric)
	if !ok {
		m.log.Warn().Interface("metric", metric).Msg("not found metric")
		w.WriteHeader(http.StatusNotFound)

		return
	}

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

func (m MetricAlerts) GetValue(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}

	metric.MType = chi.URLParam(r, "metric_type")
	metric.ID = chi.URLParam(r, "metric_name")

	err := validator.ValidateGet(metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err validate metric")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric, ok := m.metricStorage.GetMetric(metric)
	if !ok {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if metric.MType == entities.Gauge {
		w.Write([]byte(fmt.Sprintf("%v", *metric.Value)))
	} else {
		w.Write([]byte(fmt.Sprintf("%v", *metric.Delta)))
	}

	w.WriteHeader(http.StatusOK)
}
