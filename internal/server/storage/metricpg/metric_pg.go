package metricpg

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/storage"
	"metric-alert/internal/server/storage/db"
	"metric-alert/internal/server/storage/metricpg/implementation"
)

const DefaultQueryTimeout = time.Second * 15

type MetricDB interface {
	SetMetric(ctx context.Context, metric entities.Metrics) (*entities.Metrics, error)
	SetMetrics(ctx context.Context, metric []entities.Metrics) error
	GetMetric(ctx context.Context, metricID string) (*entities.Metrics, error)
	GetAllMetric(ctx context.Context) ([]entities.Metrics, error)
	GetMetricsByIDs(ctx context.Context, IDs []string) ([]entities.Metrics, error)

	Ping(ctx context.Context) error
	Close()
}

var _ MetricDB = &MetricStoragePG{}

type MetricStoragePG struct {
	implementation.MetricStorage
}

func NewMetricDB(pg *db.Postgres) *MetricStoragePG {
	return &MetricStoragePG{MetricStorage: implementation.NewMetricStorage(pg)}
}

func (m *MetricStoragePG) SetMetric(ctx context.Context, metric entities.Metrics) (*entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	tx, err := m.MetricStorage.BeginTX(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = m.MetricStorage.RollbackTX(ctx, tx); err != nil {
			logger.Log.Error().Err(err).Msg("err rollback tx")
		}
	}()

	if metric.MType == entities.Gauge {
		err = m.setMetricTX(ctx, tx, metric)
		if err != nil {
			return nil, err
		}

		err = m.CommitTX(ctx, tx)
		if err != nil {
			return nil, err
		}

		return &metric, nil
	}

	counter, err := m.MetricStorage.GetMetricByIDTX(ctx, tx, metric.ID)
	if err != nil {
		if err.Error() != storage.NotFound {
			return nil, err
		}
	} else {
		newDelta := *counter.Delta + *metric.Delta
		metric.Delta = &newDelta
	}

	err = m.setMetricTX(ctx, tx, metric)
	if err != nil {
		return nil, err
	}

	err = m.CommitTX(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m *MetricStoragePG) SetMetrics(ctx context.Context, metrics []entities.Metrics) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	tx, err := m.MetricStorage.BeginTX(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err = m.MetricStorage.RollbackTX(ctx, tx); err != nil {
			logger.Log.Error().Err(err).Msg("err rollback tx")
		}
	}()

	for _, metric := range metrics {
		if metric.MType == entities.Gauge {
			err = m.setMetricTX(ctx, tx, metric)
			if err != nil {
				return err
			}

			continue
		}

		var counter *entities.Metrics
		counter, err = m.MetricStorage.GetMetricByIDTX(ctx, tx, metric.ID)
		if err != nil {
			if err.Error() != storage.NotFound {
				return err
			}
		} else {
			newDelta := *counter.Delta + *metric.Delta
			metric.Delta = &newDelta
		}

		err = m.setMetricTX(ctx, tx, metric)
		if err != nil {
			return err
		}
	}

	err = m.CommitTX(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricStoragePG) GetMetric(ctx context.Context, metricID string) (*entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetMetricByID(ctx, metricID)
}

func (m *MetricStoragePG) GetMetricsByIDs(ctx context.Context, IDs []string) ([]entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetMetricsByIDs(ctx, IDs)
}

func (m *MetricStoragePG) GetAllMetric(ctx context.Context) ([]entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetAllMetric(ctx)
}

func (m *MetricStoragePG) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.Ping(ctx)
}

func (m *MetricStoragePG) Close() {
	m.MetricStorage.Close()
}

func (m *MetricStoragePG) setMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error {
	err := m.MetricStorage.UpsertMetricTX(ctx, tx, metric)
	if err != nil {
		return err
	}

	return nil
}
