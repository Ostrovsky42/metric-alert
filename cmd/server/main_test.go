package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"metric-alert/internal/entities"
	"metric-alert/internal/handlers"
	"metric-alert/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricValid(t *testing.T) {
	mockStorage, _ := repository.InitRepo("", "", 40, false)

	var value float64 = 5
	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
		Value: &value,
	}
	data, _ := json.Marshal(testMetric)

	req := httptest.NewRequest("POST", "/update/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil) //todo mock
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, err := mockStorage.GetMetric(testMetric.ID)
	assert.Equal(t, testMetric, finelyMetric)
	assert.Equal(t, nil, err)
}

func TestUpdateMetricInvalid(t *testing.T) {
	mockStorage, _ := repository.InitRepo("", "", 40, false)

	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
	}

	data, _ := json.Marshal(testMetric)
	req := httptest.NewRequest("POST", "/update/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	finelyMetric, err := mockStorage.GetMetric(testMetric.ID)
	assert.Equal(t, finelyMetric, entities.Metrics{})
	assert.Equal(t, errors.New("not found metric"), err)
}

func TestUpdateMetricsValid(t *testing.T) {
	mockStorage, _ := repository.InitRepo("", "", 40, false)

	var value float64 = 5
	var valueCount int64 = 5
	testMetric := []entities.Metrics{
		{
			MType: entities.Gauge,
			ID:    "test_gauge",
			Value: &value,
		},
		{
			MType: entities.Counter,
			ID:    "test_counter",
			Delta: &valueCount,
		},
	}
	data, _ := json.Marshal(testMetric)

	req := httptest.NewRequest("POST", "/updates/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil) //todo mock
	r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, err := mockStorage.GetMetricsByIDs([]string{testMetric[0].ID, testMetric[1].ID})
	assert.Equal(t, testMetric, finelyMetric)
	assert.Equal(t, nil, err)
}

func TestUpdateMetricsInvalid(t *testing.T) {
	mockStorage, _ := repository.InitRepo("", "", 40, false)

	var value float64 = 5
	testMetric := []entities.Metrics{
		{
			MType: entities.Gauge,
			ID:    "test_gauge",
			Value: &value,
		},
		{
			MType: entities.Counter,
			ID:    "test_counter",
		},
	}
	data, _ := json.Marshal(testMetric)

	req := httptest.NewRequest("POST", "/updates/", bytes.NewReader(data))

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, nil) //todo mock
	r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
