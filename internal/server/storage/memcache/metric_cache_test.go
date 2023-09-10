package memcache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"metric-alert/internal/server/entities"
)

func floatPointer(val float64) *float64 {
	return &val
}

func intPointer(val int64) *int64 {
	return &val
}

var errorNoFound = errors.New("not found metric")

func TestMemStorage_GetMetric(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemCache
		metric  entities.Metrics
		want    *entities.Metrics
		err     error
	}{
		{
			name: "positive test",
			storage: &MemCache{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Gauge, Value: floatPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric", MType: entities.Gauge},
			want:   &entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(0)},
			err:    nil,
		},
		{
			name: "positive test",
			storage: &MemCache{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric", MType: entities.Counter},
			want:   &entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(0)},
			err:    nil,
		},
		{
			name: "negative test",
			storage: &MemCache{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric-alert", MType: entities.Counter},
			want:   nil,
			err:    errorNoFound,
		},

		{
			name: "negative test",
			storage: &MemCache{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric-alert", MType: entities.Gauge},
			want:   nil,
			err:    errorNoFound,
		},
		{
			name: "positive empty name key?",
			storage: &MemCache{
				storage: map[string]entities.Metrics{"": {ID: "", MType: entities.Gauge, Value: floatPointer(8)}},
			},
			metric: entities.Metrics{ID: "", MType: entities.Gauge},
			want:   &entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(8)},
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.storage.GetMetric(context.Background(), tt.metric.ID)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.err)
		})
	}
}

func TestMemStorage_SetMetric(t *testing.T) {
	tests := []struct {
		name     string
		storage  *MemCache
		metric   entities.Metrics
		metricID string
		want     *entities.Metrics
		err      error
	}{
		{
			name:     "positive test",
			storage:  &MemCache{storage: map[string]entities.Metrics{}},
			metric:   entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(44)},
			metricID: "metric",
			want:     &entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(44)},
			err:      nil,
		},
		{
			name:     "positive test",
			storage:  &MemCache{storage: map[string]entities.Metrics{}},
			metric:   entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			metricID: "metric",
			want:     &entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			err:      nil,
		},
		{
			name: "positive test iter Counter",
			storage: func() *MemCache {
				mem := &MemCache{storage: map[string]entities.Metrics{}}
				mem.storage["metric"] = entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)}
				return mem
			}(),
			metric:   entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			metricID: "metric",
			want:     &entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(110)},
			err:      nil,
		},
		{
			name:     "negative test",
			storage:  &MemCache{storage: map[string]entities.Metrics{}},
			metric:   entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			metricID: "ne metric",
			want:     nil,
			err:      errorNoFound,
		},
		{
			name:     "negative test",
			storage:  &MemCache{storage: map[string]entities.Metrics{}},
			metric:   entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(155)},
			metricID: "ne metric",
			want:     nil,
			err:      errorNoFound,
		},
		{
			name:     "positive empty name key?",
			storage:  &MemCache{storage: map[string]entities.Metrics{}},
			metric:   entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(1)},
			metricID: "",
			want:     &entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(1)},
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.storage.SetMetric(context.Background(), tt.metric)
			if err != nil {
				log.Fatal(err)
			}

			got, got1 := tt.storage.GetMetric(context.Background(), tt.metricID)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.err)
		})
	}
}

func BenchmarkSetMetric(b *testing.B) {
	mc := NewMemCache()
	ctx := context.Background()
	metric := entities.Metrics{
		ID:    fmt.Sprintf("example_metric_%d", 1),
		MType: entities.Gauge,
		Delta: new(int64),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mc.SetMetric(ctx, metric)
		if err != nil {
			b.Fatalf("Error setting metric: %s", err)
		}
	}
}

func BenchmarkGetMetric(b *testing.B) {
	mc := NewMemCache()
	ctx := context.Background()

	metric := entities.Metrics{
		ID:    "example_metric",
		MType: entities.Gauge,
		Delta: new(int64),
	}

	_, err := mc.SetMetric(ctx, metric)
	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mc.GetMetric(ctx, metric.ID)
		if err != nil {
			b.Fatalf("Error getting metric: %s", err)
		}
	}
}

func TestMemCache_GetAllMetric(t *testing.T) {
	metrics := []entities.Metrics{
		{
			ID:    "example_metric_one",
			MType: entities.Gauge,
			Delta: new(int64),
		}, {
			ID:    "example_metric_two",
			MType: entities.Gauge,
			Delta: new(int64),
		},
	}
	tests := []struct {
		name          string
		cashedMetrics []entities.Metrics
		want          []entities.Metrics
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name:          "Test OK",
			cashedMetrics: metrics,
			want:          metrics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemCache()
			err := m.SetMetrics(context.Background(), tt.cashedMetrics)
			if err != nil {
				log.Fatal(err)
			}
			got, err := m.GetAllMetric(context.Background())
			assert.Equalf(t, tt.want, got, "GetAllMetric()")
		})
	}
}
