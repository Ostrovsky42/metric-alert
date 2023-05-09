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

func (m *MemStorage) GetMetric(metric types.Metric) (types.Metric, bool) {
	var ok bool
	if metric.MetricType == types.Gauge {
		metric.GaugeValue, ok = m.gauge[metric.MetricName]
	} else {
		metric.CounterValue, ok = m.counter[metric.MetricName]
	}

	return metric, ok
}
