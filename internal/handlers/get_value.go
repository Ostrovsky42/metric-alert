package handlers

import (
	"net/http"

	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/helpers"
)

func (m MetricAlerts) GetValue(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := helpers.UnmarshalBody(r.Body, &metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = validator.ValidateGet(metric)
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
