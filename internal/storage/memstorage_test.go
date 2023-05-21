package storage

import (
	"reflect"
	"testing"

	"metric-alert/internal/entities"
)

func floatPointer(val float64) *float64 {
	return &val
}

func intPointer(val int64) *int64 {
	return &val
}

func TestMemStorage_GetMetric(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		metric  entities.Metrics
		want    entities.Metrics
		ok      bool
	}{
		{
			name: "positive test",
			storage: &MemStorage{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Gauge, Value: floatPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric", MType: entities.Gauge},
			want:   entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(0)},
			ok:     true,
		},
		{
			name: "positive test",
			storage: &MemStorage{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric", MType: entities.Counter},
			want:   entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(0)},
			ok:     true,
		},
		{
			name: "negative test",
			storage: &MemStorage{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric-alert", MType: entities.Counter},
			want:   entities.Metrics{},
			ok:     false,
		},

		{
			name: "negative test",
			storage: &MemStorage{
				storage: map[string]entities.Metrics{"metric": {ID: "metric", MType: entities.Counter, Delta: intPointer(0)}},
			},
			metric: entities.Metrics{ID: "metric-alert", MType: entities.Gauge},
			want:   entities.Metrics{},
			ok:     false,
		},
		{
			name: "positive empty name key?",
			storage: &MemStorage{
				storage: map[string]entities.Metrics{"": {ID: "", MType: entities.Gauge, Value: floatPointer(8)}},
			},
			metric: entities.Metrics{ID: "", MType: entities.Gauge},
			want:   entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(8)},
			ok:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.storage.GetMetric(tt.metric)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("GetMetric() got1 = %v, want %v", got1, tt.ok)
			}
		})
	}
}

func TestMemStorage_SetMetric(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		metric  entities.Metrics
		want    entities.Metrics
		ok      bool
	}{
		{
			name:    "positive test",
			storage: &MemStorage{storage: map[string]entities.Metrics{}},
			metric:  entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(44)},
			want:    entities.Metrics{ID: "metric", MType: entities.Gauge, Value: floatPointer(44)},
			ok:      true,
		},
		{
			name:    "positive test",
			storage: &MemStorage{storage: map[string]entities.Metrics{}},
			metric:  entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			want:    entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			ok:      true,
		},
		{
			name:    "negative test",
			storage: &MemStorage{storage: map[string]entities.Metrics{}},
			metric:  entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(55)},
			want:    entities.Metrics{},
			ok:      false,
		},
		{
			name:    "negative test",
			storage: &MemStorage{storage: map[string]entities.Metrics{}},
			metric:  entities.Metrics{ID: "metric", MType: entities.Counter, Delta: intPointer(155)},
			want:    entities.Metrics{},
			ok:      false,
		},
		{
			name:    "positive empty name key?",
			storage: &MemStorage{storage: map[string]entities.Metrics{}},
			metric:  entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(1)},
			want:    entities.Metrics{ID: "", MType: entities.Gauge, Value: floatPointer(1)},
			ok:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.storage.SetMetric(tt.metric)
			got, got1 := tt.storage.GetMetric(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("GetMetric() got1 = %v, want %v", got1, tt.ok)
			}

		})
	}
}
