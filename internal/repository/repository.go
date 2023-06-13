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
	MetricRepo
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
		repo.MetricRepo = repo.MetricStoragePG

		return &repo, nil
	}

	memCache := memcache.NewMemCache()
	fileStorage, err := filestorage.NewFileRecorder(fileStoragePath, memCache)
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
	return r.MetricRepo.SetMetric(metric)
}

func (r *Repository) SetMetrics(metric []entities.Metrics) error {
	return r.MetricRepo.SetMetrics(metric)
}

func (r *Repository) GetMetric(metricID string) (entities.Metrics, error) {
	return r.MetricRepo.GetMetric(metricID)
}

func (r *Repository) GetMetricsByIDs(IDs []string) ([]entities.Metrics, error) {
	return r.MetricRepo.GetMetricsByIDs(IDs)
}

func (r *Repository) GetAllMetric() ([]entities.Metrics, error) {
	return r.MetricRepo.GetAllMetric()
}

func (r *Repository) Ping() error {
	return r.MetricRepo.Ping()
}

func (r *Repository) Close() {
	if r.MetricStoragePG != nil {
		r.MetricStoragePG.Close()
	}
}
