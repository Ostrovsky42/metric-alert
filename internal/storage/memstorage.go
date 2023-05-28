package storage

import (
	"metric-alert/internal/entities"
)

type MemStorage struct {
	storage map[string]entities.Metrics
}

type MetricStorage interface {
	SetMetric(metric entities.Metrics) entities.Metrics
	GetMetric(metric entities.Metrics) (entities.Metrics, bool)
	GetAllMetric() []entities.Metrics
}

func NewMemStore() *MemStorage {
	return &MemStorage{storage: make(map[string]entities.Metrics)}
}

func (m *MemStorage) SetMetric(metric entities.Metrics) entities.Metrics {
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
	return metric
}

func (m *MemStorage) GetMetric(metric entities.Metrics) (entities.Metrics, bool) {
	metric, ok := m.storage[metric.ID]
	if ok {
		return metric, ok
	}

	return metric, ok
}

func (m *MemStorage) GetAllMetric() []entities.Metrics {
	metrics := make([]entities.Metrics, 0, len(m.storage))

	for _, id := range MetricIDs {
		if metric, ok := m.storage[id]; ok {
			metrics = append(metrics, metric)
		}
	}

	return metrics
}

func (m *MemStorage) SetMetrics(metrics []entities.Metrics) {
	for _, metric := range metrics {
		m.storage[metric.ID] = metric
	}
}
