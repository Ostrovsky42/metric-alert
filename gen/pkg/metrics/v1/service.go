package v1

import (
	"context"

	"github.com/bufbuild/protovalidate-go"

	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ MetricsServiceServer = &Service{}

type Service struct {
	repo      repository.MetricRepo
	validator *protovalidate.Validator
}

func NewService(repo repository.MetricRepo) Service {
	v, err := protovalidate.New(protovalidate.WithMessages(&UpdateMetricsReq{}))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to initialize validator")
	}
	return Service{
		repo:      repo,
		validator: v,
	}
}

func (m *Service) UpdateMetricsV1(ctx context.Context, req *UpdateMetricsReq) (*emptypb.Empty, error) {
	if err := m.validator.Validate(req); err != nil {
		logger.Log.Error().Err(err).Msg("failed validation")
	}

	metrics := make([]entities.Metrics, 0, len(req.Metrics))
	var metric entities.Metrics
	for _, metricReq := range req.Metrics {
		switch metricReq.Type {
		case MetricType_GAUGE:
			value := metricReq.Value
			metric = entities.Metrics{
				ID:    metricReq.Id,
				MType: entities.Gauge,
				Value: &value,
			}

		case MetricType_COUNTER:
			delta := metricReq.Delta
			metric = entities.Metrics{
				ID:    metricReq.Id,
				MType: entities.Counter,
				Delta: &delta,
			}
		}

		metrics = append(metrics, metric)
	}

	if err := m.repo.SetMetrics(ctx, metrics); err != nil {
		logger.Log.Error().Err(err).Msg("error set metrics")
		return nil, status.Errorf(codes.Internal, "error set metrics: %s", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (m *Service) mustEmbedUnimplementedMetricsServiceServer() {}
