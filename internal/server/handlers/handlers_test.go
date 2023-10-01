package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"metric-alert/internal/server/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/repository"
	"metric-alert/internal/server/storage"
)

func TestUpdateMetricWithBody(t *testing.T) {
	type testCase struct {
		name       string
		req        *http.Request
		statusCode int
	}

	var value float64 = 5
	testCases := []testCase{
		{
			name: "ok",
			req: func() *http.Request {
				metric := entities.Metrics{
					MType: entities.Gauge,
					ID:    "test_gauge",
					Value: &value,
				}
				data, _ := json.Marshal(metric)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusOK,
		},
		{
			name: "Not found - empty metric name",
			req: func() *http.Request {
				metric := entities.Metrics{
					MType: entities.Gauge,
					ID:    "",
					Value: &value,
				}
				data, _ := json.Marshal(metric)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad request - empty metric type",
			req: func() *http.Request {
				metric := entities.Metrics{
					MType: "",
					ID:    "test",
					Value: &value,
				}
				data, _ := json.Marshal(metric)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad request - empty value",
			req: func() *http.Request {
				metric := entities.Metrics{
					MType: entities.Counter,
					ID:    "test",
				}
				data, _ := json.Marshal(metric)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad request - invalid body",
			req: func() *http.Request {
				data, _ := json.Marshal("{")
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
	}

	mockStorage, _ := repository.InitRepo("", "", 40, false)

	for _, tc := range testCases {
		rr := httptest.NewRecorder()

		r := chi.NewRouter()
		metricAlerts := NewMetric(mockStorage, nil)
		r.Post("/update/", metricAlerts.UpdateMetricWithBody)

		r.ServeHTTP(rr, tc.req)

		assert.Equal(t, tc.statusCode, rr.Code)
	}
}

func TestUpdateMetricsWithBody(t *testing.T) {
	type testCase struct {
		name       string
		req        *http.Request
		statusCode int
	}

	var value float64 = 5
	var сount int64 = 5
	testCases := []testCase{
		{
			name: "OK",
			req: func() *http.Request {
				metrics := []entities.Metrics{
					{
						MType: entities.Gauge,
						ID:    "test_gauge1",
						Value: &value,
					},
					{
						MType: entities.Counter,
						ID:    "test_counter1",
						Delta: &сount,
					},
				}
				data, _ := json.Marshal(metrics)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusOK,
		},
		{
			name: "Not found - empty metric name",
			req: func() *http.Request {
				metrics := []entities.Metrics{
					{
						MType: entities.Gauge,
						ID:    "test_gauge2",
						Value: &value,
					},
					{
						MType: entities.Counter,
						ID:    "",
						Value: &value,
					},
				}
				data, _ := json.Marshal(metrics)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusNotFound,
		},
		{
			name: "Not found - empty metric name",
			req: func() *http.Request {
				metrics := []entities.Metrics{
					{
						MType: entities.Gauge,
						ID:    "test_gauge2",
						Value: &value,
					},
					{
						MType: "",
						ID:    "",
						Value: &value,
					},
				}
				data, _ := json.Marshal(metrics)
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad request - invalid body",
			req: func() *http.Request {
				data, _ := json.Marshal("{")
				return httptest.NewRequest("POST", "/update/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
	}

	mockStorage, _ := repository.InitRepo("", "", 40, false)

	for _, tc := range testCases {
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		metricAlerts := NewMetric(mockStorage, nil)
		r.Post("/update/", metricAlerts.UpdateMetricsWithBody)

		r.ServeHTTP(rr, tc.req)

		assert.Equal(t, tc.statusCode, rr.Code)
	}
}

func TestUpdateMetric(t *testing.T) {
	type testCase struct {
		name       string
		req        *http.Request
		statusCode int
	}

	testCases := []testCase{
		{
			name:       "OK",
			req:        httptest.NewRequest("POST", "/update/gauge/test_metric/5", nil),
			statusCode: http.StatusOK,
		},
		{
			name:       "Bad request - invalid metric type",
			req:        httptest.NewRequest("POST", "/update//test_metric/5", nil),
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Not found - empty metric name",
			req:        httptest.NewRequest("POST", "/update/gauge//5", nil),
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Bad request - invalid metric value",
			req:        httptest.NewRequest("POST", "/update/gauge/test_metric/invalid_value", nil),
			statusCode: http.StatusBadRequest,
		},
	}

	mockStorage, _ := repository.InitRepo("", "", 40, false)

	for _, tc := range testCases {
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		metricAlerts := NewMetric(mockStorage, nil)
		r.Post("/update/{type}/{name}/{value}", metricAlerts.UpdateMetric)

		r.ServeHTTP(rr, tc.req)

		assert.Equal(t, tc.statusCode, rr.Code)
	}
}

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

	metricAlerts := NewMetric(mockStorage, nil)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, err := mockStorage.GetMetric(context.Background(), testMetric.ID)
	assert.Equal(t, &testMetric, finelyMetric)
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

	metricAlerts := NewMetric(mockStorage, nil)
	r.Post("/update/", metricAlerts.UpdateMetricWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	_, err := mockStorage.GetMetric(context.Background(), testMetric.ID)
	assert.Equal(t, errors.New(storage.NotFound), err)
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

	metricAlerts := NewMetric(mockStorage, nil)
	r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	finelyMetric, err := mockStorage.GetMetricsByIDs(context.Background(), []string{testMetric[0].ID, testMetric[1].ID})
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

	metricAlerts := NewMetric(mockStorage, nil)
	r.Post("/updates/", metricAlerts.UpdateMetricsWithBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetValueWithBody(t *testing.T) {
	type testCase struct {
		name         string
		req          *http.Request
		storedMetric *entities.Metrics
		statusCode   int
	}

	testCases := []testCase{
		{
			name: "OK",
			storedMetric: &entities.Metrics{
				ID:    "test",
				MType: entities.Counter,
			},
			req: func() *http.Request {
				testMetric := entities.Metrics{
					ID:    "test",
					MType: entities.Counter,
				}
				data, _ := json.Marshal(testMetric)
				return httptest.NewRequest("POST", "/get_value/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusOK,
		},
		{
			name: "Bad request - empty type",
			req: func() *http.Request {
				data, _ := json.Marshal(entities.Metrics{ID: "non_existent_metric"})
				return httptest.NewRequest("POST", "/get_value/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad request - invalid body",
			req: func() *http.Request {
				data, _ := json.Marshal("{")
				return httptest.NewRequest("POST", "/get_value/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not found - metric not found",
			req: func() *http.Request {
				data, _ := json.Marshal(entities.Metrics{ID: "non_existent_metric", MType: entities.Gauge})
				return httptest.NewRequest("POST", "/get_value/", bytes.NewReader(data))
			}(),
			statusCode: http.StatusNotFound,
		},
	}

	mockStorage, _ := repository.InitRepo("", "", 40, false)

	for _, tc := range testCases {
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		if tc.storedMetric != nil {
			if _, err := mockStorage.SetMetric(context.Background(), *tc.storedMetric); err != nil {
				logger.Log.Fatal().Err(err).Send()
			}
		}
		metricAlerts := NewMetric(mockStorage, nil)
		r.Post("/get_value/", metricAlerts.GetValueWithBody)

		r.ServeHTTP(rr, tc.req)

		assert.Equal(t, tc.statusCode, rr.Code)
	}
}
