package storage

import (
	"metric-alert/internal/server/entities"
	"reflect"
	"testing"
)

var sampleMetrics = []entities.Metrics{
	{ID: "id1", MType: entities.Gauge},
	{ID: "id2", MType: entities.Counter},
	{ID: "id5", MType: entities.Gauge},
	{ID: "id3", MType: entities.Gauge},
	{ID: "id4", MType: entities.Counter},
}

var sampleIDs = []string{"id1", "id2", "id3", "id1", "id2", "id4"}

func BenchmarkSortMetric(b *testing.B) {
	metricsToSort := make([]entities.Metrics, len(sampleMetrics))
	copy(metricsToSort, sampleMetrics)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SortMetric(metricsToSort)
	}
}

func BenchmarkRemoveDuplicatesIDs(b *testing.B) {
	idsToRemoveDuplicates := make([]string, len(sampleIDs))
	copy(idsToRemoveDuplicates, sampleIDs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RemoveDuplicatesIDs(idsToRemoveDuplicates)
	}
}

func TestSortMetric(t *testing.T) {
	tests := []struct {
		name    string
		metrics []entities.Metrics
		want    []entities.Metrics
	}{
		{
			name: "Test sort 1",
			metrics: []entities.Metrics{
				{ID: "Poll", MType: "Counter"},
				{ID: "Allo", MType: "Gauge"},
				{ID: "Alloc", MType: "Gauge"},
			},
			want: []entities.Metrics{
				{ID: "Allo", MType: "Gauge"},
				{ID: "Alloc", MType: "Gauge"},
				{ID: "Poll", MType: "Counter"},
			},
		},
		{
			name: "Test sort 2",
			metrics: []entities.Metrics{
				{ID: "Allo", MType: "Gauge"},
				{ID: "Poll", MType: "Counter"},
				{ID: "Alloc", MType: "Gauge"},
			},
			want: []entities.Metrics{
				{ID: "Allo", MType: "Gauge"},
				{ID: "Alloc", MType: "Gauge"},
				{ID: "Poll", MType: "Counter"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortMetric(tt.metrics)
			if !reflect.DeepEqual(tt.metrics, tt.want) {
				t.Errorf("SortMetric() = %v, want %v", tt.metrics, tt.want)
			}
		})
	}
}
