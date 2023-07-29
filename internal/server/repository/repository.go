package repository

import (
	"context"
	"time"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/storage/db"
	"metric-alert/internal/server/storage/filestorage"
	"metric-alert/internal/server/storage/memcache"
	"metric-alert/internal/server/storage/metricpg"
)

type MetricRepo interface {
	SetMetric(ctx context.Context, metric entities.Metrics) (*entities.Metrics, error)
	SetMetrics(ctx context.Context, metric []entities.Metrics) error
	GetMetric(ctx context.Context, metricID string) (*entities.Metrics, error)
	GetAllMetric(ctx context.Context) ([]entities.Metrics, error)
	GetMetricsByIDs(ctx context.Context, IDs []string) ([]entities.Metrics, error)

	Ping(ctx context.Context) error
	Close()
}

type Repository struct {
	*memcache.MemCache
	*metricpg.MetricStoragePG
	*filestorage.FileRecorder
}

func InitRepo(fileStoragePath, dataBaseDSN string, storeIntervalSec int, restore bool) (MetricRepo, error) {
	var repo Repository
	if dataBaseDSN != "" {
		pg, err := db.NewPostgresDB(dataBaseDSN)
		if err != nil {
			return nil, err
		}
		repo.MetricStoragePG = metricpg.NewMetricDB(pg)

		return repo.MetricStoragePG, nil
	}

	repo.MemCache = memcache.NewMemCache()

	fileStorage, err := filestorage.NewFileRecorder(fileStoragePath, repo.MemCache)
	if err != nil {
		return nil, err
	}
	repo.FileRecorder = fileStorage

	if restore {
		fileStorage.RestoreMetrics()
	}

	go repo.StartRecording(storeIntervalSec)

	return repo.MemCache, nil
}

func (r *Repository) StartRecording(updateInterval int) {
	interval := time.Duration(updateInterval) * time.Second

	for {
		time.Sleep(interval)
		r.FileRecorder.RecordMetrics()
	}
}

func (r *Repository) Close() {
	if r.MetricStoragePG != nil {
		r.MetricStoragePG.Close()
	}
}
