package memcache

import (
	"context"
	"errors"
	"fmt"

	"metric-alert/internal/entities"
	"metric-alert/internal/storage"
)

type MemCache struct {
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

func (m *MemCache) SetMetric(ctx context.Context, metric entities.Metrics) (*entities.Metrics, error) {
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

func (m *MemCache) GetMetric(ctx context.Context, metricID string) (*entities.Metrics, error) {
	metric, ok := m.storage[metricID]
	if ok {
		return &metric, nil
	}

	return nil, errors.New(storage.NotFound)
}

func (m *MemCache) GetMetricsByIDs(ctx context.Context, IDs []string) ([]entities.Metrics, error) {
	var metrics []entities.Metrics
	IDs = storage.RemoveDuplicatesIDs(IDs)
	for _, id := range IDs {
		metric, ok := m.storage[id]
		if !ok {
			return nil, fmt.Errorf("%s by id:%s", storage.NotFound, id)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *MemCache) GetAllMetric(ctx context.Context) ([]entities.Metrics, error) {
	metrics := make([]entities.Metrics, 0, len(m.storage))

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

func (m *MemCache) Ping(ctx context.Context) error {
	return nil
}

func (m *MemCache) Close() {
}
