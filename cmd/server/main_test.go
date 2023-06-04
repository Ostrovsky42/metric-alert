package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"metric-alert/internal/entities"
	"metric-alert/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"metric-alert/internal/storage"
)

func TestUpdateMetricValid(t *testing.T) {
	mockStorage := storage.NewMemStore()

	var value float64 = 5
	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
		Value: &value,
	}
	data, _ := json.Marshal(testMetric)

	req := httptest.NewRequest("POST", "/update/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil, nil) //todo mock
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric.ID)
	assert.Equal(t, testMetric, finelyMetric)
	assert.Equal(t, true, ok)
}

func TestUpdateMetricInvalid(t *testing.T) {
	mockStorage := storage.NewMemStore()

	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
	}

	data, _ := json.Marshal(testMetric)
	req := httptest.NewRequest("POST", "/update/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil, nil)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric.ID)
	assert.Equal(t, finelyMetric, entities.Metrics{})
	assert.Equal(t, false, ok)
}
