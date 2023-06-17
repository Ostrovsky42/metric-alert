package agent

import (
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
}

func NewAgent(reportInterval, pollInterval int, serverURL, signKey string) *Agent {
	return &Agent{
		sender:         metricsender.NewMetricSender(serverURL, signKey),
		gatherer:       gatherer.NewGatherer(pollInterval),
		reportInterval: time.Duration(reportInterval) * time.Second,
	}
}

func (a *Agent) Run() {
	go a.gatherer.GatherMetrics()
	//go a.sendReport()
	a.sendPackReportJSON()
}

func (a *Agent) sendPackReportJSON() {
	for {
		var metrics []gatherer.Metrics
		time.Sleep(a.reportInterval)
		for _, metric := range a.gatherer.Metrics {
			metrics = append(metrics, metric)
		}
		if len(metrics) == 0 {
			continue
		}

		if err := a.sender.SendMetricPackJSON(metrics); err != nil {
			logger.Log.Error().Err(err).Msg("err SendMetricPackJSON")
		}
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
