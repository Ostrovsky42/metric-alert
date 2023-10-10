package main

import (
	"context"
	"errors"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "metric-alert/gen/pkg/metrics/v1"
	"metric-alert/internal/server/config"
	"metric-alert/internal/server/handlers"
	"metric-alert/internal/server/logger"
	"metric-alert/internal/server/repository"
)

const templatePath = "internal/server/html/templates/info_page.html"
const shutdownTimeout = 10

type Application struct {
	metric         handlers.MetricAlerts
	storage        repository.MetricRepo
	isHTTP         bool
	serverHost     string
	signKey        string
	privateKeyPath string
	subnet         string
}

func NewApp(cfg *config.Config) Application {
	memRepo, err := repository.InitRepo(
		cfg.FileStoragePath,
		cfg.DataBaseDSN,
		cfg.StoreIntervalSec,
		cfg.Restore,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed init storage")
	}

	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error while parse web templates")
	}

	return Application{
		metric:         handlers.NewMetric(memRepo, tmp),
		storage:        memRepo,
		serverHost:     cfg.ServerHost,
		signKey:        cfg.SignKey,
		privateKeyPath: cfg.CryptoKey,
		subnet:         cfg.TrustedSubnet,
		isHTTP:         cfg.IsHTTP,
	}
}

func (a Application) Run() {
	switch a.isHTTP {
	case true:
		a.RunHTTP()
	case false:
		a.RunGRPC()
	}
}

func (a Application) RunHTTP() {
	s := http.Server{
		Addr:    a.serverHost,
		Handler: NewRoutes(a.metric, a.signKey, a.privateKeyPath, a.subnet),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal().Err(err).Msg("Error start serve")
		}
	}()

	<-a.shutdownSignal()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Error Shutdown server")
	}
}

func (a Application) RunGRPC() {
	grpcServer := grpc.NewServer()

	service := pb.NewService(a.storage)
	pb.RegisterMetricsServiceServer(grpcServer, &service)
	lis, err := net.Listen("tcp", a.serverHost)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error start listen")
	}

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Log.Fatal().Err(err).Msg("Error start serve")
		}
	}()

	<-a.shutdownSignal()

	grpcServer.GracefulStop()
}

func (a Application) Close() {
	a.storage.Close()
}

func (a Application) shutdownSignal() chan os.Signal {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	return shutdownSignal
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
