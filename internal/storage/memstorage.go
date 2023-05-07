package storage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type MetricStorage interface {
	SetGauge(name string, gauge float64)
	Count(name string, value int64)
}

func NewMemStore() MetricStorage {
	g := make(map[string]float64)
	c := make(map[string]int64)
	return &MemStorage{counter: c, gauge: g}
}

func (m *MemStorage) SetGauge(name string, gauge float64) {
	m.gauge[name] = gauge
}

func (m *MemStorage) Count(name string, value int64) {
	m.counter[name] += value
}
