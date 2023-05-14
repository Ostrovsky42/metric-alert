package main

import (
	"github.com/go-chi/chi"
	"metric-alert/internal/handlers"
	"metric-alert/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"metric-alert/internal/storage"
)

func TestUpdateMetricValid(t *testing.T) {
	mockStorage := storage.NewMemStore()

	testMetric := types.Metric{
		MetricType: types.Gauge,
		MetricName: "test_gauge",
		GaugeValue: 5,
	}

	req := httptest.NewRequest("POST", "/update/gauge/test_gauge/5", nil)

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage)
	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricAlerts.UpdateMetric)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric)
	assert.Equal(t, testMetric, finelyMetric)
	assert.Equal(t, true, ok)
}

func TestUpdateMetricInvalid(t *testing.T) {
	mockStorage := storage.NewMemStore()

	testMetric := types.Metric{
		MetricType: types.Gauge,
		MetricName: "test_gauge",
		GaugeValue: 6,
	}

	req := httptest.NewRequest("POST", "/update/counter/test_counter/5", nil)

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage)
	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricAlerts.UpdateMetric)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric)
	assert.NotEqual(t, testMetric, finelyMetric)
	assert.Equal(t, false, ok)
}
