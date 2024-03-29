package metricsender

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"metric-alert/internal/agent/compressor"
	"metric-alert/internal/agent/gatherer"
	"metric-alert/internal/crypto/hybrid"
	"metric-alert/internal/crypto/symmetric"
	"metric-alert/internal/hasher"
	"metric-alert/internal/server/logger"
)

const numberOfAttempts = 3

type MetricSender struct {
	client            *http.Client
	hashBuilder       hasher.HashBuilder
	encryptor         *hybrid.Encryptor
	serverURL         string
	attemptsIntervals []int
}

func NewMetricSender(serverURL string, signKey string, cryptoKeyPath string) *MetricSender {
	encryptor, err := hybrid.NewEncryptor(cryptoKeyPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err create encryptor")
	}

	return &MetricSender{
		client:            &http.Client{},
		hashBuilder:       hasher.NewHashGenerator(signKey),
		encryptor:         encryptor,
		serverURL:         "http://" + serverURL,
		attemptsIntervals: []int{1, 3, 5},
	}
}

func (s *MetricSender) SendMetricPackJSON(metrics []gatherer.Metrics) error {
	if len(metrics) == 0 {
		logger.Log.Info().Msg("empty metrics")

		return nil
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("json.Marshal :%w", err)
	}
	metricURL := fmt.Sprintf("%s/updates/", s.serverURL)

	var iv []byte
	if s.encryptor.Included() {
		data, iv, err = s.encryptor.Symmetric.Encrypt(data)
		if err != nil {
			return fmt.Errorf("encryptor.Encrypt :%w", err)
		}
	}
	compressed, err := compressor.CompressData(data)
	if err != nil {
		return fmt.Errorf("compressor.CompressData :%w", err)
	}
	req, err := http.NewRequest("POST", metricURL, compressed)
	if err != nil {
		return fmt.Errorf("http.NewRequest :%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	s.signRequest(data, iv, req)

	var resp *http.Response
	for i := 0; i < numberOfAttempts; i++ {
		resp, err = s.client.Do(req)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) ||
				errors.Is(err, syscall.ECONNREFUSED) {
				logger.Log.Warn().Interface("req", req).Err(err).Int("attempt", i+1).
					Msg("unsuccessful attempt send request")

				time.Sleep(time.Duration(s.attemptsIntervals[i]) * time.Second)

				continue
			}

			return fmt.Errorf("client.Do :%w", err)
		}

		if err = resp.Body.Close(); err != nil {
			return fmt.Errorf("resp.Body.Close :%w", err)
		}

		break
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

func (s *MetricSender) signRequest(data, iv []byte, req *http.Request) {
	if s.encryptor.Included() {
		req.Header.Set(symmetric.HeaderSymmetricKey, s.encryptor.EncryptedKey)
		req.Header.Set(symmetric.HeaderInitializationVector, base64.StdEncoding.EncodeToString(iv))
	}
	if s.hashBuilder.IsNotActive() {
		return
	}

	hash := s.hashBuilder.GetHash(data)
	req.Header.Set("HashSHA256", hash)
}
