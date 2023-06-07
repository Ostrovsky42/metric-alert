package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"

	"metric-alert/internal/logger"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgresDB(dsn string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{DB: pool}

	if err = pg.initTable(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close() {
	logger.Log.Info().Msg("CLOSE DB")
	p.DB.Close()
}

func (p *Postgres) initTable() error {
	ctx := context.Background()

	_, err := p.DB.Exec(ctx, `CREATE TYPE metrics_type AS ENUM ('counter', 'gauge')`)
	if err == nil {
		logger.Log.Info().Msg("success created type metrics_type")
	}

	table := `
CREATE TABLE IF NOT EXISTS metrics
(
    id varchar unique not null,
    metric_type metrics_type not null,
    value double precision,
    delta bigint
);`
	_, err = p.DB.Exec(ctx, table)
	if err != nil {
		logger.Log.Info().Msg("err created table metrics")

		return err
	}

	return nil
}
