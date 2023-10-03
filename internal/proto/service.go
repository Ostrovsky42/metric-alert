package proto

import (
	"context"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ MetricsServiceServer = &Service{}

type Service struct {
	repo repository.MetricRepo
}

func NewService(repo repository.MetricRepo) Service {
	return Service{
		repo: repo,
	}
}

func (m *Service) UpdateMetrics(ctx context.Context, req *UpdateMetricsReq) (*emptypb.Empty, error) {
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
