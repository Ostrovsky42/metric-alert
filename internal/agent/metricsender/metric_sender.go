package metricsender

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "metric-alert/gen/pkg/metrics/v1"
	"metric-alert/internal/agent/compressor"
	"metric-alert/internal/agent/gatherer"
	"metric-alert/internal/crypto/hybrid"
	"metric-alert/internal/crypto/symmetric"
	"metric-alert/internal/hasher"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
)

const numberOfAttempts = 3

type MetricSender struct {
	httpClient        *http.Client
	grpcClient        pb.MetricsServiceClient
	hashBuilder       hasher.HashBuilder
	encryptor         *hybrid.Encryptor
	serverURL         string
	localIP           string
	isHTTP            bool
	attemptsIntervals []int
}

func NewMetricSender(isHTTP bool, serverURL string, localIP string, signKey string, cryptoKeyPath string) *MetricSender {
	encryptor, err := hybrid.NewEncryptor(cryptoKeyPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err create encryptor")
	}
	ms := &MetricSender{
		httpClient:        &http.Client{},
		hashBuilder:       hasher.NewHashGenerator(signKey),
		encryptor:         encryptor,
		serverURL:         "http://" + serverURL,
		localIP:           localIP,
		attemptsIntervals: []int{1, 3, 5},
	}
	if isHTTP {
		return ms
	}

	conn, err := grpc.Dial(serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err dial conn")
	}
	ms.grpcClient = pb.NewMetricsServiceClient(conn)

	return ms
}

func (s *MetricSender) SendMetricPack(metrics []gatherer.Metrics) error {
	if s.isHTTP {
		return s.sendMetricPackJSON(metrics)
	}
	return s.sendMetricPackGRPC(metrics)
}

func (s *MetricSender) sendMetricPackJSON(metrics []gatherer.Metrics) error {
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
	req.Header.Set("X-Real-IP", s.localIP)
	s.signRequest(data, iv, req)

	var resp *http.Response
	for i := 0; i < numberOfAttempts; i++ {
		resp, err = s.httpClient.Do(req)
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

func (s *MetricSender) sendMetricPackGRPC(metrics []gatherer.Metrics) error {
	grpcMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, metric := range metrics {
		grpcMetric := &pb.Metric{
			Id: metric.ID,
		}

		switch metric.MType {
		case entities.Gauge:
			grpcMetric.Type = pb.MetricType_GAUGE
			grpcMetric.Value = getValue(metric.Value)
		case entities.Counter:
			grpcMetric.Type = pb.MetricType_COUNTER
			grpcMetric.Delta = metric.Delta
		}

		grpcMetrics = append(grpcMetrics, grpcMetric)
	}

	md := metadata.New(map[string]string{"X-Real-IP": s.localIP})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := s.grpcClient.UpdateMetricsV1(ctx, &pb.UpdateMetricsReq{
		Metrics: grpcMetrics,
	}); err != nil {
		return fmt.Errorf("send metrics process error: %w", err)
	}

	return nil
}

func getValue(value any) float64 {
	switch v := value.(type) {
	case uint64:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	}

	return 0
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
	response, err := s.httpClient.Do(req)
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
	response, err := s.httpClient.Do(req)
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
