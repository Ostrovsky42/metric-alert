package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"metric-alert/internal/server/logger"
)

func TrustedSubnetInterceptor(subnet string) grpc.UnaryServerInterceptor {
	if len(subnet) == 0 {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed parse CIDR")
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		ipr := md.Get("X-Real-IP")
		if len(ipr) != 1 {
			logger.Log.Info().Interface("invalid_ip", ipr).Msg("invalid number of ip")
			return nil, status.Error(codes.PermissionDenied, "Invalid IP")
		}

		ip := net.ParseIP(ipr[0])
		if !ipNet.Contains(ip) {
			logger.Log.Info().Str("invalid_ip", ip.String()).Send()
			return nil, status.Error(codes.PermissionDenied, "Invalid IP")
		}

		return handler(ctx, req)
	}
}
