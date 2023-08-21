package filestorage

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/storage/memcache"
)

const perm = 0666

type FileRecorder struct {
	filename    string
	metricCache *memcache.MemCache
}

func NewFileRecorder(
	filename string,
	memStorage *memcache.MemCache,
) *FileRecorder {
	return &FileRecorder{
		filename:    filename,
		metricCache: memStorage,
	}
}

func (f *FileRecorder) RestoreMetrics() {
	file, err := os.OpenFile(f.filename, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err open file")
	}
	defer file.Close()

	var metrics []entities.Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil && err != io.EOF {
		logger.Log.Error().Err(err).Msg("err file Decoder")
	}

	f.metricCache.SetMetrics(context.Background(), metrics)
}

func (f *FileRecorder) RecordMetrics() {
	metrics, _ := f.metricCache.GetAllMetric(context.Background())
	if len(metrics) == 0 {
		return
	}

	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err open file")

		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(metrics)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err update file")
	}
}
