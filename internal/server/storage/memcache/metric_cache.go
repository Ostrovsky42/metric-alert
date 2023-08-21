package memcache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/storage"
)

type MemCache struct {
	mu      sync.RWMutex
	storage map[string]entities.Metrics
}

type MetricCache interface {
	SetMetric(ctx context.Context, metric entities.Metrics) (*entities.Metrics, error)
	SetMetrics(ctx context.Context, metric []entities.Metrics) error
	GetMetric(ctx context.Context, metricID string) (*entities.Metrics, error)
	GetAllMetric(ctx context.Context) ([]entities.Metrics, error)
	GetMetricsByIDs(ctx context.Context, IDs []string) ([]entities.Metrics, error)

	Ping(ctx context.Context) error
	Close()
}

var _ MetricCache = &MemCache{}

func NewMemCache() *MemCache {
	return &MemCache{storage: make(map[string]entities.Metrics)}
}

func (m *MemCache) SetMetric(_ context.Context, metric entities.Metrics) (*entities.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	return &metric, nil
}

func (m *MemCache) GetMetric(_ context.Context, metricID string) (*entities.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metric, ok := m.storage[metricID]
	if ok {
		return &metric, nil
	}

	return nil, errors.New(storage.NotFound)
}

func (m *MemCache) GetMetricsByIDs(_ context.Context, ids []string) ([]entities.Metrics, error) {
	ids = storage.RemoveDuplicatesIDs(ids)
	metrics := make([]entities.Metrics, 0, len(ids))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, id := range ids {
		metric, ok := m.storage[id]
		if !ok {
			return nil, fmt.Errorf("%s by id:%s", storage.NotFound, id)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *MemCache) GetAllMetric(_ context.Context) ([]entities.Metrics, error) {
	metrics := make([]entities.Metrics, 0, len(m.storage))
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, metric := range m.storage {
		metrics = append(metrics, metric)
	}

	storage.SortMetric(metrics)

	return metrics, nil
}

func (m *MemCache) SetMetrics(ctx context.Context, metrics []entities.Metrics) error {
	for _, metric := range metrics {
		if _, err := m.SetMetric(ctx, metric); err != nil {
			return err
		}
	}

	return nil
}

func (m *MemCache) Ping(_ context.Context) error {
	return nil
}

func (m *MemCache) Close() {
}
