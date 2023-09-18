// Package handlers предоставляет обработчики HTTP-запросов для работы с метриками.
package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/handlers/validator"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/repository"
)

// @Title MetricAlerts API
// @Description Metric storage service.

// @host      localhost:8080

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name authorization

// @Tag.name GetMetric
// @Tag.description "Group of requests to get service metrics"

// @Tag.name UpdateMetric
// @Tag.description "Group of requests to update service metrics"

const (
	metricType  = "type"
	metricName  = "name"
	metricValue = "value"
)

// MetricAlerts представляет структуру для обработчиков метрик.
type MetricAlerts struct {
	metricStorage repository.MetricRepo
	tmp           *template.Template
}

// NewMetric создает новый экземпляр MetricAlerts.
func NewMetric(metricStorage repository.MetricRepo, tmp *template.Template) MetricAlerts {
	return MetricAlerts{
		metricStorage: metricStorage,
		tmp:           tmp,
	}
}

// UpdateMetricWithBody предоставляет обработчик для обновления метрики из тела запроса.
// UpdateMetricWithBody godoc
// @Summary Update metric
// @Description  Update metric from request body
// @Tags         UpdateMetric
// @Accept       json
// @Produce      json
// @Param		 metric_data body entities.Metrics true "Metric data"
// @Success      200  {object}  entities.Metrics
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /update/ [post].
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
		if err.Error() == validator.EmptyMetricName {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	updatedMetric, err := m.metricStorage.SetMetric(r.Context(), metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error set metric")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(updatedMetric)
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

// UpdateMetricsWithBody предоставляет обработчик для обновления метрик из тела запроса.
// UpdateMetricsWithBody godoc
// @Summary Update metrics
// @Description  Update metrics from request body
// @Tags         UpdateMetric
// @Accept       json
// @Produce      json
// @Param		 metrics_data body []entities.Metrics true "Arrays metric data"
// @Success      200  {object}  []entities.Metrics
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /updates/ [post].
func (m MetricAlerts) UpdateMetricsWithBody(w http.ResponseWriter, r *http.Request) {
	var metrics []entities.Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err unmarshal body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	err = validator.ValidateMetrics(metrics)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metrics")
		if err.Error() == validator.EmptyMetricName {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	ctx := r.Context()
	err = m.metricStorage.SetMetrics(ctx, metrics)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error set metrics")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	updatedMetricIDs := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		updatedMetricIDs = append(updatedMetricIDs, metric.ID)
	}

	updatedMetric, err := m.metricStorage.GetMetricsByIDs(ctx, updatedMetricIDs)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error get metrics by ids")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = json.Marshal(updatedMetric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err encode data")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateMetric предоставляет обработчик для обновления метрики по типу, имени и значению из пути.
// UpdateMetric godoc
// @Summary Update metrics
// @Description  Update metric by specifying its type, name, and value from path
// @Tags         UpdateMetric
// @Accept       json
// @Produce      json
// @Param        type   path      string  true  "Metric Type"
// @Param        name   path      string  true  "Metric Name"
// @Param        value  path      number  true  "Metric Value"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /update/{type}/{name}/{value}/ [post].
func (m MetricAlerts) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := entities.Metrics{}

	metric.MType = chi.URLParam(r, metricType)
	metric.ID = chi.URLParam(r, metricName)
	mValue := chi.URLParam(r, metricValue)

	err := validator.ValidateUpdate(&metric, mValue)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err validate metric")
		if err.Error() == validator.EmptyMetricName {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = m.metricStorage.SetMetric(r.Context(), metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error set metric")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
