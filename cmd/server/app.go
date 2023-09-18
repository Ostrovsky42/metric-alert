package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	serverHost     string
	signKey        string
	privateKeyPath string
}

func NewApp(cfg config.Config) Application {
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
	}
}

func (a Application) Run() {
	s := http.Server{
		Addr:    a.serverHost,
		Handler: NewRoutes(a.metric, a.signKey, a.privateKeyPath),
	}

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal().Err(err).Msg("Error start serve")
		}
	}()

	<-shutdownSignal

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Error Shutdown server")
	}
}

func (a Application) Close() {
	a.storage.Close()
}

/*
_ = template.FuncMap{
"floatPoint": func(p *float64) float64 { return *p },
}   <td>{{floatPoint .Value | printf "%.0f"}}</td>
*/
