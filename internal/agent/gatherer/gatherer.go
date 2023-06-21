package gatherer

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"metric-alert/internal/server/entities"
	"metric-alert/internal/server/logger"
)

type Gatherer struct {
	Metrics      map[string]Metrics
	pollInterval time.Duration
	mu           sync.RWMutex
}

func NewGatherer(pollInterval int) *Gatherer {
	return &Gatherer{
		Metrics:      make(map[string]Metrics, 0),
		pollInterval: time.Duration(pollInterval) * time.Second,
	}
}

func (g *Gatherer) GatherRuntimeMetrics() {
	var m runtime.MemStats
	var delta int64
	runtime.ReadMemStats(&m)
	g.mu.Lock()

	g.Metrics[Alloc] = Metrics{ID: Alloc, MType: entities.Gauge, Value: m.Alloc}
	g.Metrics[BuckHashSys] = Metrics{ID: BuckHashSys, MType: entities.Gauge, Value: m.BuckHashSys}
	g.Metrics[Frees] = Metrics{ID: Frees, MType: entities.Gauge, Value: m.Frees}
	g.Metrics[GCSys] = Metrics{ID: GCSys, MType: entities.Gauge, Value: m.GCSys}
	g.Metrics[HeapAlloc] = Metrics{ID: HeapAlloc, MType: entities.Gauge, Value: m.HeapAlloc}
	g.Metrics[HeapIdle] = Metrics{ID: HeapIdle, MType: entities.Gauge, Value: m.HeapIdle}
	g.Metrics[HeapInuse] = Metrics{ID: HeapInuse, MType: entities.Gauge, Value: m.HeapInuse}
	g.Metrics[HeapObjects] = Metrics{ID: HeapObjects, MType: entities.Gauge, Value: m.HeapObjects}
	g.Metrics[HeapReleased] = Metrics{ID: HeapReleased, MType: entities.Gauge, Value: m.HeapReleased}
	g.Metrics[HeapSys] = Metrics{ID: HeapSys, MType: entities.Gauge, Value: m.HeapSys}
	g.Metrics[LastGC] = Metrics{ID: LastGC, MType: entities.Gauge, Value: m.LastGC}
	g.Metrics[Lookups] = Metrics{ID: Lookups, MType: entities.Gauge, Value: m.Lookups}
	g.Metrics[MCacheInuse] = Metrics{ID: MCacheInuse, MType: entities.Gauge, Value: m.MCacheInuse}
	g.Metrics[MCacheSys] = Metrics{ID: MCacheSys, MType: entities.Gauge, Value: m.MCacheSys}
	g.Metrics[MSpanInuse] = Metrics{ID: MSpanInuse, MType: entities.Gauge, Value: m.MSpanInuse}
	g.Metrics[MSpanSys] = Metrics{ID: MSpanSys, MType: entities.Gauge, Value: m.MSpanSys}
	g.Metrics[Mallocs] = Metrics{ID: Mallocs, MType: entities.Gauge, Value: m.Mallocs}
	g.Metrics[NextGC] = Metrics{ID: NextGC, MType: entities.Gauge, Value: m.NextGC}
	g.Metrics[NumForcedGC] = Metrics{ID: NumForcedGC, MType: entities.Gauge, Value: m.NumForcedGC}
	g.Metrics[NumGC] = Metrics{ID: NumGC, MType: entities.Gauge, Value: m.NumGC}
	g.Metrics[OtherSys] = Metrics{ID: OtherSys, MType: entities.Gauge, Value: m.OtherSys}
	g.Metrics[PauseTotalNs] = Metrics{ID: PauseTotalNs, MType: entities.Gauge, Value: m.PauseTotalNs}
	g.Metrics[StackInuse] = Metrics{ID: StackInuse, MType: entities.Gauge, Value: m.StackInuse}
	g.Metrics[StackSys] = Metrics{ID: StackSys, MType: entities.Gauge, Value: m.StackSys}
	g.Metrics[Sys] = Metrics{ID: Sys, MType: entities.Gauge, Value: m.Sys}
	g.Metrics[TotalAlloc] = Metrics{ID: TotalAlloc, MType: entities.Gauge, Value: m.TotalAlloc}
	g.Metrics[GCCPUFraction] = Metrics{ID: GCCPUFraction, MType: entities.Gauge, Value: m.GCCPUFraction}
	g.Metrics[RandomValue] = Metrics{ID: RandomValue, MType: entities.Gauge, Value: rand.Uint32()}
	g.Metrics[PollCount] = Metrics{ID: PollCount, MType: entities.Counter, Delta: delta}

	g.mu.Unlock()
	delta++
}

func (g *Gatherer) GatherMemoryMetrics() {
	v, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error().Err(err).Msg("err gatherer virtual memory metrics")

		return
	}

	cpuNum, err := cpu.Counts(false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("err gatherer cpu metrics")

		return
	}

	g.mu.Lock()

	g.Metrics[TotalMemory] = Metrics{ID: TotalMemory, MType: entities.Gauge, Value: v.Total}
	g.Metrics[FreeMemory] = Metrics{ID: TotalMemory, MType: entities.Gauge, Value: v.Free}
	g.Metrics[CPUutilization1] = Metrics{ID: CPUutilization1, MType: entities.Gauge, Value: int64(cpuNum)}

	g.mu.Unlock()
}

func (g *Gatherer) GetMetricToSend() []Metrics {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var metrics []Metrics
	for _, metric := range g.Metrics {
		metrics = append(metrics, metric)
	}

	return metrics
}

func (g *Gatherer) StartMetricsGatherer() {
	ticker := time.NewTicker(g.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go g.GatherRuntimeMetrics()
			go g.GatherMemoryMetrics()
		}
	}
}
