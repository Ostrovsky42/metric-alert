package storage

import (
	"encoding/json"
	"io"
	"metric-alert/internal/entities"
	"metric-alert/internal/logger"
	"os"
)

type FileRecorder struct {
	filename    string
	metricCache *MemCache
	isRestore   bool
}

func NewFileRecorder(
	filename string,
	memStorage *MemCache,
) (*FileRecorder, error) {
	return &FileRecorder{
		filename:    filename,
		metricCache: memStorage,
	}, nil
}

func (f *FileRecorder) RestoreMetrics() {
	file, err := os.OpenFile(f.filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err open file")
	}
	defer file.Close()

	var metrics []entities.Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil && err != io.EOF {
		logger.Log.Error().Err(err).Msg("err file Decoder")
	}

	f.metricCache.SetMetrics(metrics)
}

func (f *FileRecorder) RecordMetrics() {
	metrics, _ := f.metricCache.GetAllMetric()
	if len(metrics) == 0 {
		return
	}

	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
