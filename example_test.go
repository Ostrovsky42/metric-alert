package metric_alert_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/handlers"
	"metric-alert/internal/server/repository"
)

func ExampleMetricAlerts_UpdateMetricWithBody() {
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
}

func ExampleMetricAlerts_UpdateMetricsWithBody() {
	mockStorage, err := repository.InitRepo("", "", 40, false)
	if err != nil {
		log.Fatal("err init repo: ", err.Error())
	}

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

	metricAlerts := handlers.NewMetric(mockStorage, nil)
	r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
}

func ExampleMetricAlerts_GetValueWithBody() {
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
}
