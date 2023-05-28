package storage

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"io"
	"metric-alert/internal/entities"
	"os"
	"time"
)

type FileRecorder struct {
	file           *os.File
	memStorage     *MemStorage
	updateInterval time.Duration
	isRestore      bool
	log            zerolog.Logger
}

func NewFileRecorder(filename string, interval int, restore bool, memStorage *MemStorage, log zerolog.Logger) (*FileRecorder, error) {
	openParam := os.O_RDWR | os.O_CREATE
	if !restore {
		openParam |= os.O_TRUNC
	}

	file, err := os.OpenFile(filename, openParam, 0666)
	if err != nil {
		return nil, err
	}

	return &FileRecorder{
		file:           file,
		memStorage:     memStorage,
		isRestore:      restore,
		updateInterval: time.Duration(interval) * time.Second,
		log:            log,
	}, nil
}

func (f *FileRecorder) Run() {
	defer f.file.Close()

	f.restore()

	f.recordMetric()
}

func (f *FileRecorder) restore() {
	if !f.isRestore {
		return
	}

	var metrics []entities.Metrics
	err := json.NewDecoder(f.file).Decode(&metrics)
	if err != nil && err != io.EOF {
		f.log.Error().Err(err).Msg("err file Decoder")
	}

	f.memStorage.SetMetrics(metrics)
}

func (f *FileRecorder) recordMetric() {
	e := json.NewEncoder(f.file)
	for {
		time.Sleep(f.updateInterval)
		metrics := f.memStorage.GetAllMetric()
		if len(metrics) > 0 {
			err := f.clearFile()
			if err != nil {
				f.log.Error().Err(err).Msg("err clear file")

				continue
			}

			err = e.Encode(metrics)
			if err != nil {
				f.log.Error().Err(err).Msg("err update file")
			}
		}
	}
}

func (f *FileRecorder) clearFile() error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
