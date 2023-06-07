package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"metric-alert/internal/entities"
	"metric-alert/internal/storage/postgres"
	"time"
)

const DefaultQueryTimeout = time.Second * 5

var _ MetricStorage = &MetricPG{}

type MetricPG struct {
	*postgres.Postgres
}

func NewMetricDB(pg *postgres.Postgres) *MetricPG {
	return &MetricPG{Postgres: pg}
}

func (m *MetricPG) SetMetric(metric entities.Metrics) (entities.Metrics, error) {
	if metric.MType == entities.Gauge {
		err := m.setMetric(metric)
		if err != nil {
			return entities.Metrics{}, err
		}

		return metric, nil
	}

	counter, err := m.GetMetric(metric.ID)
	if err != nil {
		return entities.Metrics{}, err
	}

	newDelta := *counter.Delta + *metric.Delta
	metric.Delta = &newDelta

	err = m.setMetric(metric)
	if err != nil {
		return entities.Metrics{}, err
	}

	return metric, nil
}

func (m *MetricPG) setMetric(metric entities.Metrics) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	updateQuery := `UPDATE metrics SET value = $1, delta = $2 WHERE id = $3`
	result, err := m.DB.Exec(ctxTimeout, updateQuery, metric.Value, metric.Delta, metric.ID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected != 0 {
		return nil
	}

	insertQuery := `INSERT INTO metrics (id, metric_type, value, delta)  VALUES ($1, $2, $3, $4)`
	_, err = m.DB.Query(ctxTimeout, insertQuery,
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

func (m *MetricPG) GetMetric(metricID string) (entities.Metrics, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	sql := `SELECT id, metric_type, value, delta FROM metrics WHERE id = $1;`
	var metric entities.Metrics

	err := m.DB.QueryRow(ctxTimeout, sql, metricID).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Value,
		&metric.Delta,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Metrics{}, errors.New("not found metric")
		}

		return entities.Metrics{}, err
	}

	return metric, nil
}

func (m *MetricPG) GetAllMetric() ([]entities.Metrics, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()

	sql := `SELECT id, metric_type, value, delta FROM metrics;`
	rows, err := m.DB.Query(ctxTimeout, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entities.Metrics
	for rows.Next() {
		var metric entities.Metrics
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	sortMetric(metrics)

	return metrics, nil
}

func (m *MetricPG) SetMetrics(metrics []entities.Metrics) {
	//	sql := `INSERT INTO metrics (id, metric_type, value, delta) VALUES `
	//values := make([]interface{}, 0, len(m.cache.storage)*4)

	for _, metric := range metrics {
		m.setMetric(metric)

		//sql += fmt.Sprintf("($%d, $%d, $%d, $%d),", i*4+1, i*4+2, i*4+3, i*4+4)
		//values = append(values, metric.ID, metric.MType, metric.Value, metric.Delta)
	}
	// Удаление последней запятой
	//sql = sql[:len(sql)-1]

	//_, err := m.DB.Exec(ctx, sql, values...)
	//if err != nil {
	//	// Обработка ошибки
	//	return
	//}
}

func (m *MetricPG) Ping() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	return m.DB.Ping(ctxTimeout)
}
