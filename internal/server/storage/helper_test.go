package storage

import (
	"metric-alert/internal/server/entities"
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
