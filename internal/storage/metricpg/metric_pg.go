package metricpg

import (
	"context"
	"github.com/jackc/pgx/v4"
	"metric-alert/internal/entities"
	"metric-alert/internal/logger"
	"metric-alert/internal/storage"
	"metric-alert/internal/storage/db"
	"metric-alert/internal/storage/metricpg/implementation"
	"time"
)

const DefaultQueryTimeout = time.Second * 15

type MetricDB interface {
	SetMetric(metric entities.Metrics) (entities.Metrics, error)
	SetMetrics(metric []entities.Metrics) error
	GetMetric(metricID string) (entities.Metrics, error)
	GetAllMetric() ([]entities.Metrics, error)
	GetMetricsByIDs(IDs []string) ([]entities.Metrics, error)

	Ping() error
}

var _ MetricDB = &MetricStoragePG{}

type MetricStoragePG struct {
	implementation.MetricStorage
}

func NewMetricDB(pg *db.Postgres) *MetricStoragePG {
	return &MetricStoragePG{MetricStorage: implementation.NewMetricStorage(pg)}
}

func (m *MetricStoragePG) SetMetric(metric entities.Metrics) (entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	tx, err := m.MetricStorage.BeginTX(ctx)
	if err != nil {
		return entities.Metrics{}, err
	}

	defer func() {
		if err = m.MetricStorage.RollbackTX(ctx, tx); err != nil {
			logger.Log.Error().Err(err).Msg("err rollback tx")
		}
	}()

	if metric.MType == entities.Gauge {
		err = m.setMetricTX(ctx, tx, metric)
		if err != nil {
			return entities.Metrics{}, err
		}

		err = m.CommitTX(ctx, tx)
		if err != nil {
			return entities.Metrics{}, err
		}

		return metric, nil
	}

	counter, err := m.MetricStorage.GetMetricByIDTX(ctx, tx, metric.ID)
	if err != nil {
		if err.Error() != storage.NotFound {
			return entities.Metrics{}, err
		}
	} else {
		newDelta := *counter.Delta + *metric.Delta
		metric.Delta = &newDelta
	}

	err = m.setMetricTX(ctx, tx, metric)
	if err != nil {
		return entities.Metrics{}, err
	}

	err = m.CommitTX(ctx, tx)
	if err != nil {
		return entities.Metrics{}, err
	}

	return metric, nil
}

func (m *MetricStoragePG) SetMetrics(metrics []entities.Metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
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

		var counter entities.Metrics
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

func (m *MetricStoragePG) GetMetric(metricID string) (entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetMetricByID(ctx, metricID)
}

func (m *MetricStoragePG) GetMetricsByIDs(IDs []string) ([]entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetMetricsByIDs(ctx, IDs)
}

func (m *MetricStoragePG) GetAllMetric() ([]entities.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.GetAllMetric(ctx)
}

func (m *MetricStoragePG) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	return m.MetricStorage.Ping(ctx)
}

func (m *MetricStoragePG) setMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error {
	err := m.MetricStorage.UpsertMetricTX(ctx, tx, metric)
	if err != nil {
		return err
	}

	return nil
}
