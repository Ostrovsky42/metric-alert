package repository

import (
	"time"

	"metric-alert/internal/entities"
	"metric-alert/internal/storage/db"
	"metric-alert/internal/storage/filestorage"
	"metric-alert/internal/storage/memcache"
	"metric-alert/internal/storage/metricpg"
)

type MetricRepo interface {
	SetMetric(metric entities.Metrics) (entities.Metrics, error)
	SetMetrics(metric []entities.Metrics) error
	GetMetric(metricID string) (entities.Metrics, error)
	GetAllMetric() ([]entities.Metrics, error)
	GetMetricsByIDs(IDs []string) ([]entities.Metrics, error)

	Ping() error
}

var _ MetricRepo = &Repository{}

type Repository struct {
	*memcache.MemCache
	*metricpg.MetricStoragePG
	*filestorage.FileRecorder
}

func InitRepo(fileStoragePath, dataBaseDSN string, storeIntervalSec int, restore bool) (*Repository, error) {
	var repo Repository
	if dataBaseDSN != "" {
		pg, err := db.NewPostgresDB(dataBaseDSN)
		if err != nil {
			return nil, err
		}
		repo.MetricStoragePG = metricpg.NewMetricDB(pg)

		return &repo, nil
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

	return &repo, nil
}

func (r *Repository) StartRecording(updateInterval int) {
	interval := time.Duration(updateInterval) * time.Second

	for {
		time.Sleep(interval)
		r.FileRecorder.RecordMetrics()
	}
}

func (r *Repository) SetMetric(metric entities.Metrics) (entities.Metrics, error) {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.SetMetric(metric)
	}

	return r.MemCache.SetMetric(metric)
}

func (r *Repository) SetMetrics(metric []entities.Metrics) error {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.SetMetrics(metric)
	}

	return r.MemCache.SetMetrics(metric)
}

func (r *Repository) GetMetric(metricID string) (entities.Metrics, error) {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.GetMetric(metricID)
	}

	return r.MemCache.GetMetric(metricID)
}

func (r *Repository) GetMetricsByIDs(IDs []string) ([]entities.Metrics, error) {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.GetMetricsByIDs(IDs)
	}

	return r.MemCache.GetMetricsByIDs(IDs)
}

func (r *Repository) GetAllMetric() ([]entities.Metrics, error) {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.GetAllMetric()
	}

	return r.MemCache.GetAllMetric()
}

func (r *Repository) Ping() error {
	if r.MetricStoragePG != nil {
		return r.MetricStoragePG.Ping()
	}

	return r.MemCache.Ping()
}

func (r *Repository) Close() {
	if r.MetricStoragePG != nil {
		r.MetricStoragePG.Close()
	}
}
