package metricsender

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"metric-alert/internal/agent/gatherer"
	"metric-alert/internal/compressor"
	"metric-alert/internal/logger"
)

const numberOfAttempts = 3

type MetricSender struct {
	client            *http.Client
	serverURL         string
	attemptsIntervals []int
}

func NewMetricSender(serverURL string) *MetricSender {
	return &MetricSender{
		client:            &http.Client{},
		serverURL:         "http://" + serverURL,
		attemptsIntervals: []int{1, 3, 5},
	}
}

func (s *MetricSender) SendMetricPackJSON(metrics []gatherer.Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	metricURL := fmt.Sprintf("%s/updates/", s.serverURL)
	compressed, err := compressor.CompressData(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", metricURL, compressed)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	for i := 0; i < numberOfAttempts; i++ {
		_, err = s.client.Do(req)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) {
				logger.Log.Warn().Interface("req", req).Err(err).Int("attempt", i+1).
					Msg("unsuccessful attempt send request")

				time.Sleep(time.Duration(s.attemptsIntervals[i]) * time.Second)

				continue

			}
			return err
		}
	}

	return nil
}

func (s *MetricSender) SendMetricJSON(metric gatherer.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	metricURL := fmt.Sprintf("%s/update/", s.serverURL)
	compressed, err := compressor.CompressData(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", metricURL, compressed)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	response, err := s.client.Do(req)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *MetricSender) SendMetric(mType string, name string, value interface{}) error {
	metricURL := fmt.Sprintf("%s/update/%s/%s/%v", s.serverURL, mType, name, value)
	req, err := http.NewRequest("POST", metricURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	response, err := s.client.Do(req)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
