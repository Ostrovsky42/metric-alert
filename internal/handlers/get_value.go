package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"net/http"

	"metric-alert/internal/entities"
	"metric-alert/internal/handlers/validator"
	"metric-alert/internal/logger"
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

	metric, err = m.metricStorage.GetMetric(metric.ID)
	if err != nil {
		logger.Log.Warn().Interface("metric", metric).Msg("error get metric")
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)

			return
		}
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

	metric, err = m.metricStorage.GetMetric(metric.ID)
	if err != nil {
		logger.Log.Warn().Interface("metric", metric).Msg("error get metric")
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	sendOK(w, metric)
}

func (m MetricAlerts) InfoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics, err := m.metricStorage.GetAllMetric()
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
	if err := m.metricStorage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
