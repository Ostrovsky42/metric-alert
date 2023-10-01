package middleware

import (
	"metric-alert/internal/server/logger"
	"net"
	"net/http"
)

type TrustedSubnetMiddleware struct {
	ipNet     *net.IPNet
	isInclude bool
}

func NewTrustedSubnetMW(subnet string) TrustedSubnetMiddleware {
	if len(subnet) == 0 {
		return TrustedSubnetMiddleware{}
	}

	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed parse CIDR")
	}

	return TrustedSubnetMiddleware{ipNet: ipNet, isInclude: true}
}

func (ts TrustedSubnetMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ts.isInclude {
			next.ServeHTTP(w, r)
		}
		ipr := r.Header.Get("X-Real-IP")

		ip := net.ParseIP(ipr)
		if !ts.ipNet.Contains(ip) {
			logger.Log.Info().Str("invalid_ip", ip.String()).Send()
			http.Error(w, "Invalid IP", http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
