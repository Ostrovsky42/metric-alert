package agent

import (
	"bytes"
	"compress/gzip"
	"metric-alert/internal/entities"
)

const (
	Alloc = iota
	BuckHashSys
	Frees
	GCCPUFraction
	GCSys
	HeapAlloc
	HeapIdle
	HeapInuse
	HeapObjects
	HeapReleased
	HeapSys
	LastGC
	Lookups
	MCacheInuse
	MCacheSys
	MSpanInuse
	MSpanSys
	Mallocs
	NextGC
	NumForcedGC
	NumGC
	OtherSys
	PauseTotalNs
	StackInuse
	StackSys
	Sys
	TotalAlloc
	RandomValue
	PollCount
)

func setMetricArray() *[29]entities.Metrics {
	var delta int64
	var metrics [29]entities.Metrics

	metrics[Alloc] = entities.Metrics{ID: "Alloc", MType: entities.Gauge}
	metrics[BuckHashSys] = entities.Metrics{ID: "BuckHashSys", MType: entities.Gauge}
	metrics[Frees] = entities.Metrics{ID: "Frees", MType: entities.Gauge}
	metrics[GCCPUFraction] = entities.Metrics{ID: "GCCPUFraction", MType: entities.Gauge}
	metrics[GCSys] = entities.Metrics{ID: "GCSys", MType: entities.Gauge}
	metrics[HeapAlloc] = entities.Metrics{ID: "HeapAlloc", MType: entities.Gauge}
	metrics[HeapIdle] = entities.Metrics{ID: "HeapIdle", MType: entities.Gauge}
	metrics[HeapInuse] = entities.Metrics{ID: "HeapInuse", MType: entities.Gauge}
	metrics[HeapObjects] = entities.Metrics{ID: "HeapObjects", MType: entities.Gauge}
	metrics[HeapReleased] = entities.Metrics{ID: "HeapReleased", MType: entities.Gauge}
	metrics[HeapSys] = entities.Metrics{ID: "HeapSys", MType: entities.Gauge}
	metrics[LastGC] = entities.Metrics{ID: "LastGC", MType: entities.Gauge}
	metrics[Lookups] = entities.Metrics{ID: "Lookups", MType: entities.Gauge}
	metrics[MCacheInuse] = entities.Metrics{ID: "MCacheInuse", MType: entities.Gauge}
	metrics[MCacheSys] = entities.Metrics{ID: "MCacheSys", MType: entities.Gauge}
	metrics[MSpanInuse] = entities.Metrics{ID: "MSpanInuse", MType: entities.Gauge}
	metrics[MSpanSys] = entities.Metrics{ID: "MSpanSys", MType: entities.Gauge}
	metrics[Mallocs] = entities.Metrics{ID: "Mallocs", MType: entities.Gauge}
	metrics[NextGC] = entities.Metrics{ID: "NextGC", MType: entities.Gauge}
	metrics[NumForcedGC] = entities.Metrics{ID: "NumForcedGC", MType: entities.Gauge}
	metrics[NumGC] = entities.Metrics{ID: "NumGC", MType: entities.Gauge}
	metrics[OtherSys] = entities.Metrics{ID: "OtherSys", MType: entities.Gauge}
	metrics[PauseTotalNs] = entities.Metrics{ID: "PauseTotalNs", MType: entities.Gauge}
	metrics[StackInuse] = entities.Metrics{ID: "StackInuse", MType: entities.Gauge}
	metrics[StackSys] = entities.Metrics{ID: "StackSys", MType: entities.Gauge}
	metrics[Sys] = entities.Metrics{ID: "Sys", MType: entities.Gauge}
	metrics[TotalAlloc] = entities.Metrics{ID: "TotalAlloc", MType: entities.Gauge}
	metrics[RandomValue] = entities.Metrics{ID: "RandomValue", MType: entities.Gauge}
	metrics[PollCount] = entities.Metrics{ID: "PollCount", MType: entities.Counter, Delta: &delta}

	return &metrics
}

func pointerUint64(val uint64) *float64 {
	floatVal := float64(val)
	return &floatVal
}

func pointerUint32(val uint32) *float64 {
	floatVal := float64(val)
	return &floatVal
}

func zipData(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
