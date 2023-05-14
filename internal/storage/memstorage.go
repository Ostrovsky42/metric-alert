package storage

import "metric-alert/internal/types"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type MetricStorage interface {
	SetMetric(metric types.Metric)
	GetMetric(metric types.Metric) (types.Metric, bool)
}

func NewMemStore() *MemStorage {
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

func (m *MemStorage) GetMetric(metric types.Metric) (types.Metric, bool) {
	var ok bool
	metric.GaugeValue, ok = m.gauge[metric.MetricName]
	if metric.MetricType == types.Counter {
		metric.CounterValue, ok = m.counter[metric.MetricName]
	}

	return metric, ok
}
