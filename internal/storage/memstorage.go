package storage

import (
	"metric-alert/internal/entities"
	"sort"
)

type MemStorage struct {
	storage map[string]entities.Metrics
}

type MetricStorage interface {
	SetMetric(metric entities.Metrics) entities.Metrics
	GetMetric(metricID string) (entities.Metrics, bool)
	GetAllMetric() []entities.Metrics
	SetMetrics(metrics []entities.Metrics)
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

func (m *MemStorage) GetMetric(metricID string) (entities.Metrics, bool) {
	metric, ok := m.storage[metricID]

	return metric, ok
}

func (m *MemStorage) GetAllMetric() []entities.Metrics {
	metrics := make([]entities.Metrics, 0, len(m.storage))

	for _, metric := range m.storage {
		metrics = append(metrics, metric)

	}

	sortFunc := func(i, j int) bool {
		if metrics[i].MType == entities.Counter {
			return false
		}

		if metrics[j].MType == entities.Counter {
			return true
		}

		return metrics[i].ID < metrics[j].ID
	}

	sort.Slice(metrics, sortFunc)

	return metrics
}

func (m *MemStorage) SetMetrics(metrics []entities.Metrics) {
	for _, metric := range metrics {
		m.storage[metric.ID] = metric
	}
}
