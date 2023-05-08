package storage

import "metric-alert/internal/types"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type MetricStorage interface {
	SetMetric(metric types.Metric)
}

func NewMemStore() MetricStorage {
	g := make(map[string]float64)
	c := make(map[string]int64)
	return &MemStorage{counter: c, gauge: g}
}

func (m *MemStorage) SetMetric(metric types.Metric) {
	if metric.MetricType == types.Gauge {
		m.gauge[metric.MetricName] = metric.GaugeValue
	} else {
		m.counter[metric.MetricName] += metric.CounterValue
	}
}
