package metric_alert_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/handlers"
	"metric-alert/internal/server/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleUpdateMetricWithBody() {
	_ = func(t *testing.T) {
		mockStorage, _ := repository.InitRepo("", "", 40, false)

		var value float64 = 5
		testMetric := entities.Metrics{
			MType: entities.Gauge,
			ID:    "test_gauge",
			Value: &value,
		}
		data, _ := json.Marshal(testMetric)

		req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(data))

		r := chi.NewRouter()

		metricAlerts := handlers.NewMetric(mockStorage, nil)
		r.Post("/update/", metricAlerts.UpdateMetricWithBody)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		finelyMetric, err := mockStorage.GetMetric(context.Background(), testMetric.ID)
		assert.Equal(t, &testMetric, finelyMetric)
		assert.Equal(t, nil, err)
	}
}

func ExampleUpdateMetricsWithBody() {
	_ = func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewReader(data))

		r := chi.NewRouter()

		metricAlerts := handlers.NewMetric(mockStorage, nil)
		r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
}

func ExampleGetValueWithBody() {
	_ = func(t *testing.T) {
		mockStorage, _ := repository.InitRepo("", "", 40, false)

		var value float64 = 5
		testMetric := entities.Metrics{
			MType: entities.Gauge,
			ID:    "test_gauge",
			Value: &value,
		}

		mockStorage.SetMetric(context.Background(), testMetric)
		data, _ := json.Marshal(testMetric)

		req := httptest.NewRequest(http.MethodGet, "/value/", bytes.NewReader(data))

		r := chi.NewRouter()

		metricAlerts := handlers.NewMetric(mockStorage, nil)
		r.Post("/value/", metricAlerts.GetValueWithBody)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		finelyMetric, err := mockStorage.GetMetric(context.Background(), testMetric.ID)
		assert.Equal(t, &testMetric, finelyMetric)
		assert.Equal(t, nil, err)
	}
}
