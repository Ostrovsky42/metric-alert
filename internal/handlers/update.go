package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/helpers"
	"metric-alert/internal/storage"
)

type MetricAlerts struct {
	metricStorage storage.MetricStorage
	log           zerolog.Logger
}

func NewMetric(metricStorage storage.MetricStorage, log zerolog.Logger) MetricAlerts {
	return MetricAlerts{
		metricStorage: metricStorage,
		log:           log,
	}
}

func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}
	err := helpers.UnmarshalBody(r.Body, &metric)
	if err != nil {
		m.log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	err = validator.ValidateUpdate(metric)
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
