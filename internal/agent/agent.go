package agent

import (
	"metric-alert/internal/agent/gatherer"
	"metric-alert/internal/agent/metricsender"
	"metric-alert/internal/logger"
	"time"

	"metric-alert/internal/entities"
)

type Agent struct {
	sender         *metricsender.MetricSender
	gatherer       *gatherer.Gatherer
	reportInterval time.Duration
}

func NewAgent(reportInterval, pollInterval int, serverURL string) *Agent {
	return &Agent{
		sender:         metricsender.NewMetricSender(serverURL),
		gatherer:       gatherer.NewGatherer(pollInterval),
		reportInterval: time.Duration(reportInterval) * time.Second,
	}
}

func (a *Agent) Run() {
	go a.gatherer.GatherMetrics()
	//go a.sendReport()
	a.sendReportJSON()
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
