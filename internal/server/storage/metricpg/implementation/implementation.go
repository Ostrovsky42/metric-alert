package implementation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/storage"
	"metric-alert/internal/server/storage/db"
)

type MetricStorage interface {
	BeginTX(ctx context.Context) (pgx.Tx, error)
	CommitTX(ctx context.Context, tx pgx.Tx) error
	RollbackTX(ctx context.Context, tx pgx.Tx) error

	UpsertMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error
	InsertMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error
	UpdateMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error
	GetMetricByIDTX(ctx context.Context, tx pgx.Tx, metricID string) (*entities.Metrics, error)
	GetMetricByID(ctx context.Context, metricID string) (*entities.Metrics, error)
	GetMetricsByIDs(ctx context.Context, metricIDs []string) ([]entities.Metrics, error)
	GetAllMetric(ctx context.Context) ([]entities.Metrics, error)

	Ping(ctx context.Context) error
	Close()
}

var _ MetricStorage = &MetricPG{}

type MetricPG struct {
	*db.Postgres
}

func NewMetricStorage(pg *db.Postgres) *MetricPG {
	return &MetricPG{Postgres: pg}
}

func (m *MetricPG) BeginTX(ctx context.Context) (pgx.Tx, error) {
	tx, err := m.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return tx, nil
}
func (m *MetricPG) RollbackTX(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return err
	}

	return nil
}

func (m *MetricPG) CommitTX(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			logger.Log.Error().Err(rollbackErr).Msg("err rollback tx")

		}
		return err
	}

	return nil
}

func (m *MetricPG) InsertMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error {
	insertQuery := `INSERT INTO metrics (id, metric_type, value, delta)  VALUES ($1, $2, $3, $4)`
	_, err := tx.Query(ctx, insertQuery,
		metric.ID,
		metric.MType,
		metric.Value,
		metric.Delta,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricPG) UpdateMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error {
	updateQuery := `UPDATE metrics SET value = $1, delta = $2 WHERE id = $3`
	result, err := tx.Exec(ctx, updateQuery, metric.Value, metric.Delta, metric.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New(storage.NotFound)
	}

	return nil
}

func (m *MetricPG) UpsertMetricTX(ctx context.Context, tx pgx.Tx, metric entities.Metrics) error {
	upsertQuery := `
		INSERT INTO metrics (id, metric_type, value, delta) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET value = $3, delta = $4
	`

	_, err := tx.Exec(ctx, upsertQuery,
		metric.ID,
		metric.MType,
		metric.Value,
		metric.Delta,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricPG) GetMetricByIDTX(ctx context.Context, tx pgx.Tx, metricID string) (*entities.Metrics, error) {
	sql := `SELECT id, metric_type, value, delta FROM metrics WHERE id = $1;`
	var metric entities.Metrics

	row := tx.QueryRow(ctx, sql, metricID)
	err := row.Scan(
		&metric.ID,
		&metric.MType,
		&metric.Value,
		&metric.Delta,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New(storage.NotFound)
		}

		return nil, err
	}

	return &metric, nil
}

func (m *MetricPG) GetMetricByID(ctx context.Context, metricID string) (*entities.Metrics, error) {
	sql := `SELECT id, metric_type, value, delta FROM metrics WHERE id = $1;`
	var metric entities.Metrics

	err := m.DB.QueryRow(ctx, sql, metricID).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Value,
		&metric.Delta,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New(storage.NotFound)
		}

		return nil, err
	}

	return &metric, nil
}

func (m *MetricPG) GetMetricsByIDs(ctx context.Context, metricIDs []string) ([]entities.Metrics, error) {
	sql := `SELECT id, metric_type, value, delta FROM metrics WHERE id=ANY($1);`
	rows, err := m.DB.Query(ctx, sql, pq.Array(metricIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entities.Metrics
	for rows.Next() {
		var metric entities.Metrics
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *MetricPG) GetAllMetric(ctx context.Context) ([]entities.Metrics, error) {
	sql := `SELECT id, metric_type, value, delta FROM metrics;`
	rows, err := m.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entities.Metrics
	for rows.Next() {
		var metric entities.Metrics
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	storage.SortMetric(metrics)

	return metrics, nil
}

func (m *MetricPG) Ping(ctx context.Context) error {
	return m.DB.Ping(ctx)
}

func (m *MetricPG) Close() {
	m.DB.Close()
}
