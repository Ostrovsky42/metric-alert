package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	Conn *pgx.Conn
}

func NewPostgresDB(dsn string) (*Postgres, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Postgres{Conn: conn}, nil
}

func (p *Postgres) Close() {
	p.Conn.Close(context.Background())
}
