package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"metric-alert/internal/entities"
	"metric-alert/internal/helpers"
)

type Agent struct {
	client         *http.Client
	metrics        *[29]entities.Metrics
	reportInterval time.Duration
	pollInterval   time.Duration
	serverURL      string
	log            zerolog.Logger
}

func NewAgent(reportInterval, pollInterval int, serverURL string, log zerolog.Logger) Agent {
	client := &http.Client{}
	metrics := setMetricArray()
	return Agent{
		client:         client,
		serverURL:      "http://" + serverURL,
		metrics:        metrics,
		reportInterval: time.Duration(reportInterval) * time.Second,
		pollInterval:   time.Duration(pollInterval) * time.Second,
		log:            log,
	}
}

func (a Agent) Run() {
	go a.gatherMetrics()
	//go a.sendReport()
	a.sendReportJSON()
}

func (a Agent) sendReportJSON() {
	for {
		time.Sleep(a.reportInterval)
		for _, metric := range a.metrics {
			if err := a.sendMetricJSON(metric); err != nil {
				a.log.Error().Err(err).Msg("err sendMetricJSON")
			}
		}
	}
}

func (a Agent) sendMetricJSON(metric entities.Metrics) error {
	data, err := helpers.EncodeData(metric)
	if err != nil {
		return err
	}
	metricURL := fmt.Sprintf("%s/update", a.serverURL)
	req, err := http.NewRequest("POST", metricURL, data)
	if err != nil {
		a.log.Err(err).Bytes("data", data.Bytes()).Msg("err prepare new request")

		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := a.client.Do(req)
	if err != nil {
		a.log.Err(err).Bytes("request.data", data.Bytes()).Msg("err send request")

		return err
	}
	err = response.Body.Close()
	if err != nil {
		a.log.Err(err).Msg("err close response body")

		return err
	}

	return nil
}

func (a Agent) sendReport() {
	for {
		time.Sleep(a.reportInterval)
		for _, metric := range a.metrics {
			if metric.MType != entities.Counter {
				if err := a.sendMetric(metric.MType, metric.ID, *metric.Value); err != nil {
					a.log.Error().Err(err).Msg("err sendMetric")
				}
				continue
			}

			if err := a.sendMetric(metric.MType, metric.ID, *metric.Delta); err != nil {
				a.log.Error().Err(err).Msg("err sendMetric")
			}
		}
	}
}

func (a Agent) sendMetric(mType string, name string, value interface{}) error {
	metricURL := fmt.Sprintf("%s/update/%s/%s/%v", a.serverURL, mType, name, value)
	req, err := http.NewRequest("POST", metricURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	response, err := a.client.Do(req)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a Agent) gatherMetrics() {
	var m runtime.MemStats
	for {
		runtime.ReadMemStats(&m)
		a.metrics[Alloc].Value = pointerUint64(m.Alloc)
		a.metrics[BuckHashSys].Value = pointerUint64(m.BuckHashSys)
		a.metrics[Frees].Value = pointerUint64(m.Frees)
		a.metrics[GCSys].Value = pointerUint64(m.GCSys)
		a.metrics[HeapAlloc].Value = pointerUint64(m.HeapAlloc)
		a.metrics[HeapIdle].Value = pointerUint64(m.HeapIdle)
		a.metrics[HeapInuse].Value = pointerUint64(m.HeapInuse)
		a.metrics[HeapObjects].Value = pointerUint64(m.HeapObjects)
		a.metrics[HeapReleased].Value = pointerUint64(m.HeapReleased)
		a.metrics[HeapSys].Value = pointerUint64(m.HeapSys)
		a.metrics[LastGC].Value = pointerUint64(m.LastGC)
		a.metrics[Lookups].Value = pointerUint64(m.Lookups)
		a.metrics[MCacheInuse].Value = pointerUint64(m.MCacheInuse)
		a.metrics[MCacheSys].Value = pointerUint64(m.MCacheSys)
		a.metrics[MSpanInuse].Value = pointerUint64(m.MSpanInuse)
		a.metrics[MSpanSys].Value = pointerUint64(m.MSpanSys)
		a.metrics[Mallocs].Value = pointerUint64(m.Mallocs)
		a.metrics[NextGC].Value = pointerUint64(m.NextGC)
		a.metrics[NumForcedGC].Value = pointerUint32(m.NumForcedGC)
		a.metrics[NumGC].Value = pointerUint32(m.NumGC)
		a.metrics[OtherSys].Value = pointerUint64(m.OtherSys)
		a.metrics[PauseTotalNs].Value = pointerUint64(m.PauseTotalNs)
		a.metrics[StackInuse].Value = pointerUint64(m.StackInuse)
		a.metrics[StackSys].Value = pointerUint64(m.StackSys)
		a.metrics[Sys].Value = pointerUint64(m.Sys)
		a.metrics[TotalAlloc].Value = pointerUint64(m.TotalAlloc)
		a.metrics[GCCPUFraction].Value = &m.GCCPUFraction
		a.metrics[RandomValue].Value = pointerUint32(rand.Uint32())

		*a.metrics[PollCount].Delta++
		time.Sleep(a.pollInterval)
	}
}
