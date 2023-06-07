package storage

import (
	"metric-alert/internal/entities"
	"metric-alert/internal/storage/postgres"
	"time"
)

type MetricStorage interface {
	SetMetric(metric entities.Metrics) (entities.Metrics, error)
	GetMetric(metricID string) (entities.Metrics, error)
	GetAllMetric() ([]entities.Metrics, error)
	SetMetrics(metrics []entities.Metrics)

	Ping() error
}

var _ MetricStorage = &Storage{}

type Storage struct {
	*MemCache
	*MetricPG
	*FileRecorder
}

func InitStorage(fileStoragePath, dataBaseDSN string, storeIntervalSec int, restore bool) (*Storage, error) {
	var storage Storage
	if dataBaseDSN != "" {
		pg, err := postgres.NewPostgresDB(dataBaseDSN)
		if err != nil {
			return nil, err
		}
		storage.MetricPG = NewMetricDB(pg)

		return &storage, nil
	}

	storage.MemCache = NewMemStore()
	fileStorage, err := NewFileRecorder(fileStoragePath, storage.MemCache)
	if err != nil {
		return nil, err
	}
	storage.FileRecorder = fileStorage

	if restore {
		fileStorage.RestoreMetrics()
	}

	go storage.StartRecording(storeIntervalSec)

	return &storage, nil
}

func (s *Storage) StartRecording(updateInterval int) {
	interval := time.Duration(updateInterval) * time.Second

	for {
		time.Sleep(interval)
		s.FileRecorder.RecordMetrics()
	}
}

/*
При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d
или при их пустых значениях вернитесь последовательно к:
хранению метрик в файле при наличии соответствующей переменной окружения или флага командной строки;
хранению метрик в памяти.
*/

func (s *Storage) SetMetric(metric entities.Metrics) (entities.Metrics, error) {
	if s.MetricPG != nil {
		return s.MetricPG.SetMetric(metric)
	}

	return s.MemCache.SetMetric(metric)
}

func (s *Storage) GetMetric(metricID string) (entities.Metrics, error) {
	if s.MetricPG != nil {
		return s.MetricPG.GetMetric(metricID)
	}

	return s.MemCache.GetMetric(metricID)
}

func (s *Storage) GetAllMetric() ([]entities.Metrics, error) {
	if s.MetricPG != nil {
		return s.MetricPG.GetAllMetric()
	}

	return s.MemCache.GetAllMetric()
}

func (s *Storage) SetMetrics(metrics []entities.Metrics) {
	if s.MetricPG != nil {
		s.MetricPG.SetMetrics(metrics)
	}

	s.MemCache.SetMetrics(metrics)
}

func (s *Storage) Ping() error {
	if s.MetricPG != nil {
		return s.MetricPG.Ping()
	}

	return s.MemCache.Ping()
}

func (s *Storage) Close() {
	if s.MetricPG != nil {
		s.MetricPG.Close()
	}
}
