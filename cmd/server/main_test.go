package main

import (
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"metric-alert/internal/entities"
	"metric-alert/internal/handlers"
	"metric-alert/internal/helpers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"metric-alert/internal/storage"
)

func TestUpdateMetricValid(t *testing.T) {
	mockStorage := storage.NewMemStore()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var value float64 = 5
	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
		Value: &value,
	}
	data, _ := helpers.EncodeData(testMetric)

	req := httptest.NewRequest("POST", "/update/", data)

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, log)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric)
	assert.Equal(t, testMetric, finelyMetric)
	assert.Equal(t, true, ok)
}

func TestUpdateMetricInvalid(t *testing.T) {
	mockStorage := storage.NewMemStore()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	testMetric := entities.Metrics{
		MType: entities.Gauge,
		ID:    "test_gauge",
	}

	data, _ := helpers.EncodeData(testMetric)
	req := httptest.NewRequest("POST", "/update/", data)

	r := chi.NewRouter()

	metricAlerts := handlers.NewMetric(mockStorage, log)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	finelyMetric, ok := mockStorage.GetMetric(testMetric)
	assert.NotEqual(t, testMetric, finelyMetric)
	assert.Equal(t, false, ok)
}
