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

// GetValueWithBody предоставляет обработчик для получения метрики из тела запроса.
// GetValueWithBody godoc
// @Summary Get metric
// @Description  Get metric from request body
// @Tags         GetMetric
// @Param		 metric_data body entities.Metrics true "Metric data"
// @Success      200  {object}    entities.Metrics
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /value [get].
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
	_, err = w.Write(data)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error write to json response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetValue предоставляет обработчик для получения метрики по типу, имени и значению из пути.
// GetValue godoc
// @Summary Get metric
// @Description  Get metric value by specifying its type, name, and value from path
// @Tags         GetMetric
// @Param        type   path      string  true  "Metric Type"
// @Param        name   path      string  true  "Metric Name"
// @Param        value  path      number  true  "Metric Value"
// @Success      200  {object}    entities.Metrics
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /value/{type}/{name} [get].
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

// InfoPage предоставляет обработчик для получения HTML-страницы с информацией о метриках.
// InfoPage godoc
// @Summary Get information page
// @Description Get an HTML page with information about metrics
// @Tags InfoPage
// @Produce html
// @Success 200 {string} html "HTML content"
// @Failure 500 "Internal Server Error"
// @Router / [get].
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
