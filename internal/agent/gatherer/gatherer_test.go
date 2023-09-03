package gatherer

import "testing"

func BenchmarkGatherRuntimeMetrics(b *testing.B) {
	g := NewGatherer(0)
	var delta int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GatherRuntimeMetrics(&delta)
	}
}

func BenchmarkGatherMemoryMetrics(b *testing.B) {
	g := NewGatherer(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GatherMemoryMetrics()
	}
}
