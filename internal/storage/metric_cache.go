package storage

import (
	"errors"
	"metric-alert/internal/entities"
)

type MemCache struct {
	storage map[string]entities.Metrics
}

var _ MetricStorage = &MemCache{}

func NewMemStore() *MemCache {
	return &MemCache{storage: make(map[string]entities.Metrics)}
}

func (m *MemCache) SetMetric(metric entities.Metrics) (entities.Metrics, error) {
	if metric.MType == entities.Gauge {
		m.storage[metric.ID] = metric
	} else {
		counter, ok := m.storage[metric.ID]
		if ok {
			newDelta := *counter.Delta + *metric.Delta
			metric.Delta = &newDelta
		}

		m.storage[metric.ID] = metric
	}
	return metric, nil
}

func (m *MemCache) GetMetric(metricID string) (entities.Metrics, error) {
	metric, ok := m.storage[metricID]
	if ok {
		return metric, nil
	}

	return metric, errors.New("not found metric")
}

func (m *MemCache) GetAllMetric() ([]entities.Metrics, error) {
	metrics := make([]entities.Metrics, 0, len(m.storage))

	for _, metric := range m.storage {
		metrics = append(metrics, metric)
	}

	sortMetric(metrics)

	return metrics, nil
}

func (m *MemCache) SetMetrics(metrics []entities.Metrics) {
	for _, metric := range metrics {
		m.storage[metric.ID] = metric
	}
}

func (m *MemCache) Ping() error {
	return nil
}
