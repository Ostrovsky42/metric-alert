// Package filestorage предоставляет функциональность для записи и восстановления метрик из файлового хранилища.
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

// FileRecorder представляет структуру для сохранения и восстановления метрик в файловом хранилище.
type FileRecorder struct {
	filename    string
	metricCache *memcache.MemCache
}

// NewFileRecorder создает и возвращает новый экземпляр FileRecorder с заданными параметрами.
func NewFileRecorder(
	filename string,
	memStorage *memcache.MemCache,
) *FileRecorder {
	return &FileRecorder{
		filename:    filename,
		metricCache: memStorage,
	}
}

// RestoreMetrics восстанавливает метрики из файла и сохраняет их в кеше.
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

	err = f.metricCache.SetMetrics(context.Background(), metrics)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err set restored metrics")
	}
}

// RecordMetrics сохраняет текущие метрики в файловом хранилище.
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
