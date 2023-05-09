package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type metricType string

const count = "PollCount"

const (
	gauge   metricType = "gauge"
	counter metricType = "counter"
)

type Agent struct {
	client         *http.Client
	metrics        map[string]interface{}
	reportInterval time.Duration
	pollInterval   time.Duration
	serverURL      string
}

func NewAgent(reportInterval, pollInterval int, serverURL string) Agent {
	client := &http.Client{}
	metrics := make(map[string]interface{})
	return Agent{
		client:         client,
		serverURL:      "http://" + serverURL,
		metrics:        metrics,
		reportInterval: time.Duration(reportInterval) * time.Second,
		pollInterval:   time.Duration(pollInterval) * time.Second,
	}
}

func (a Agent) Run() {
	go a.gatherMetrics()
	a.sendReport()
}

func (a Agent) sendReport() {
	pollCount := 0
	for {
		time.Sleep(a.reportInterval)
		for name, value := range a.metrics {
			if err := a.sendMetric(gauge, name, value); err != nil {
				log.Default().Println(err)
			}
		}
		pollCount++
		if err := a.sendMetric(counter, count, pollCount); err != nil {
			log.Default().Println(err)
		}
	}
}

func (a Agent) sendMetric(mType metricType, name string, value interface{}) error {
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
		a.metrics["Alloc"] = float64(m.Alloc)
		a.metrics["BuckHashSys"] = float64(m.BuckHashSys)
		a.metrics["Frees"] = float64(m.Frees)
		a.metrics["GCCPUFraction"] = m.GCCPUFraction
		a.metrics["GCSys"] = float64(m.GCSys)
		a.metrics["HeapAlloc"] = float64(m.HeapAlloc)
		a.metrics["HeapIdle"] = float64(m.HeapIdle)
		a.metrics["HeapInuse"] = float64(m.HeapInuse)
		a.metrics["HeapObjects"] = float64(m.HeapObjects)
		a.metrics["HeapReleased"] = float64(m.HeapReleased)
		a.metrics["HeapSys"] = float64(m.HeapSys)
		a.metrics["LastGC"] = float64(m.LastGC)
		a.metrics["Lookups"] = float64(m.Lookups)
		a.metrics["MCacheInuse"] = float64(m.MCacheInuse)
		a.metrics["MCacheSys"] = float64(m.MCacheSys)
		a.metrics["MSpanInuse"] = float64(m.MSpanInuse)
		a.metrics["MSpanSys"] = float64(m.MSpanSys)
		a.metrics["Mallocs"] = float64(m.Mallocs)
		a.metrics["NextGC"] = float64(m.NextGC)
		a.metrics["NumForcedGC"] = float64(m.NumForcedGC)
		a.metrics["NumGC"] = float64(m.NumGC)
		a.metrics["OtherSys"] = float64(m.OtherSys)
		a.metrics["PauseTotalNs"] = float64(m.PauseTotalNs)
		a.metrics["StackInuse"] = float64(m.StackInuse)
		a.metrics["StackSys"] = float64(m.StackSys)
		a.metrics["Sys"] = float64(m.Sys)
		a.metrics["TotalAlloc"] = float64(m.TotalAlloc)
		a.metrics["RandomValue"] = rand.Float64() * 100
		time.Sleep(a.pollInterval)
	}
}
