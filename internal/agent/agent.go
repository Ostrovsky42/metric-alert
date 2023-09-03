package agent

import (
	"metric-alert/internal/agent/config"
	"metric-alert/internal/agent/gatherer"
	"metric-alert/internal/agent/metricsender"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
	"time"
)

type Agent struct {
	sender         *metricsender.MetricSender
	gatherer       *gatherer.Gatherer
	reportInterval time.Duration
	rateLimit      int
}

func NewAgent(cfg config.Config) *Agent {
	return &Agent{
		sender:         metricsender.NewMetricSender(cfg.ServerHost, cfg.SignKey),
		gatherer:       gatherer.NewGatherer(cfg.PollIntervalSec),
		reportInterval: time.Duration(cfg.ReportIntervalSec) * time.Second,
		rateLimit:      cfg.RateLimit,
	}
}

func (a *Agent) Run() {
	go a.gatherer.StartMetricsGatherer()

	for id := 1; id <= a.rateLimit; id++ {
		go a.sendPackReportJSON(id)
	}
}

func (a *Agent) sendPackReportJSON(workerID int) {
	for {
		time.Sleep(a.reportInterval)

		metrics := a.gatherer.GetMetricToSend()
		if len(metrics) == 0 {
			continue
		}

		if err := a.sender.SendMetricPackJSON(metrics); err != nil {
			logger.Log.Error().Err(err).Msg("err SendMetricPackJSON")
		}
		logger.Log.Info().Int("worker_id", workerID).Msg("sent metrics")
	}
}

func (a *Agent) sendReportJSON() {
	for {
		time.Sleep(a.reportInterval)
		for _, metric := range a.gatherer.Metrics {
			if err := a.sender.SendMetricJSON(metric); err != nil {
				logger.Log.Error().Err(err).Msg("err sendMetricJSON")
			}
		}
	}
}

func (a *Agent) sendReport() {
	for {
		time.Sleep(a.reportInterval)
		for _, metric := range a.gatherer.Metrics {
			if metric.MType != entities.Counter {
				if err := a.sender.SendMetric(metric.MType, metric.ID, metric.Value); err != nil {
					logger.Log.Error().Err(err).Msg("err sendMetric")
				}

				continue
			}

			if err := a.sender.SendMetric(metric.MType, metric.ID, metric.Delta); err != nil {
				logger.Log.Error().Err(err).Msg("err sendMetric")
			}
		}
	}
}
