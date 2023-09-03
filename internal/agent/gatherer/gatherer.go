// Пакет gatherer предоставляет функции для сбора метрик о работе приложения.
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

// Gatherer представляет сборщик метрик и их хранилище.
type Gatherer struct {
	Metrics      map[int]Metrics
	pollInterval time.Duration
	mu           sync.RWMutex
}

// NewGatherer создает новый экземпляр сборщика метрик с заданным интервалом сбора.
func NewGatherer(pollInterval int) *Gatherer {
	return &Gatherer{
		Metrics:      make(map[int]Metrics, DefaultMetricCount),
		pollInterval: time.Duration(pollInterval) * time.Second,
	}
}

// GatherRuntimeMetrics собирает и обновляет метрики о работе приложения в реальном времени.
func (g *Gatherer) GatherRuntimeMetrics(delta *int64) {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	g.mu.Lock()

	g.Metrics[KeyAlloc] = Metrics{ID: Alloc, MType: entities.Gauge, Value: m.Alloc}
	g.Metrics[KeyBuckHashSys] = Metrics{ID: BuckHashSys, MType: entities.Gauge, Value: m.BuckHashSys}
	g.Metrics[KeyFrees] = Metrics{ID: Frees, MType: entities.Gauge, Value: m.Frees}
	g.Metrics[KeyGCSys] = Metrics{ID: GCSys, MType: entities.Gauge, Value: m.GCSys}
	g.Metrics[KeyHeapAlloc] = Metrics{ID: HeapAlloc, MType: entities.Gauge, Value: m.HeapAlloc}
	g.Metrics[KeyHeapIdle] = Metrics{ID: HeapIdle, MType: entities.Gauge, Value: m.HeapIdle}
	g.Metrics[KeyHeapInuse] = Metrics{ID: HeapInuse, MType: entities.Gauge, Value: m.HeapInuse}
	g.Metrics[KeyHeapObjects] = Metrics{ID: HeapObjects, MType: entities.Gauge, Value: m.HeapObjects}
	g.Metrics[KeyHeapReleased] = Metrics{ID: HeapReleased, MType: entities.Gauge, Value: m.HeapReleased}
	g.Metrics[KeyHeapSys] = Metrics{ID: HeapSys, MType: entities.Gauge, Value: m.HeapSys}
	g.Metrics[KeyLastGC] = Metrics{ID: LastGC, MType: entities.Gauge, Value: m.LastGC}
	g.Metrics[KeyLookups] = Metrics{ID: Lookups, MType: entities.Gauge, Value: m.Lookups}
	g.Metrics[KeyMCacheInuse] = Metrics{ID: MCacheInuse, MType: entities.Gauge, Value: m.MCacheInuse}
	g.Metrics[KeyMCacheSys] = Metrics{ID: MCacheSys, MType: entities.Gauge, Value: m.MCacheSys}
	g.Metrics[KeyMSpanInuse] = Metrics{ID: MSpanInuse, MType: entities.Gauge, Value: m.MSpanInuse}
	g.Metrics[KeyMSpanSys] = Metrics{ID: MSpanSys, MType: entities.Gauge, Value: m.MSpanSys}
	g.Metrics[KeyMallocs] = Metrics{ID: Mallocs, MType: entities.Gauge, Value: m.Mallocs}
	g.Metrics[KeyNextGC] = Metrics{ID: NextGC, MType: entities.Gauge, Value: m.NextGC}
	g.Metrics[KeyNumForcedGC] = Metrics{ID: NumForcedGC, MType: entities.Gauge, Value: m.NumForcedGC}
	g.Metrics[KeyNumGC] = Metrics{ID: NumGC, MType: entities.Gauge, Value: m.NumGC}
	g.Metrics[KeyOtherSys] = Metrics{ID: OtherSys, MType: entities.Gauge, Value: m.OtherSys}
	g.Metrics[KeyPauseTotalNs] = Metrics{ID: PauseTotalNs, MType: entities.Gauge, Value: m.PauseTotalNs}
	g.Metrics[KeyStackInuse] = Metrics{ID: StackInuse, MType: entities.Gauge, Value: m.StackInuse}
	g.Metrics[KeyStackSys] = Metrics{ID: StackSys, MType: entities.Gauge, Value: m.StackSys}
	g.Metrics[KeySys] = Metrics{ID: Sys, MType: entities.Gauge, Value: m.Sys}
	g.Metrics[KeyTotalAlloc] = Metrics{ID: TotalAlloc, MType: entities.Gauge, Value: m.TotalAlloc}
	g.Metrics[KeyGCCPUFraction] = Metrics{ID: GCCPUFraction, MType: entities.Gauge, Value: m.GCCPUFraction}
	g.Metrics[KeyRandomValue] = Metrics{ID: RandomValue, MType: entities.Gauge, Value: rand.Uint32()}
	g.Metrics[KeyPollCount] = Metrics{ID: PollCount, MType: entities.Counter, Delta: *delta}

	g.mu.Unlock()
	*delta++
}

// GatherMemoryMetrics собирает метрики о памяти и процессоре.
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

	g.Metrics[KeyTotalMemory] = Metrics{ID: TotalMemory, MType: entities.Gauge, Value: v.Total}
	g.Metrics[KeyFreeMemory] = Metrics{ID: FreeMemory, MType: entities.Gauge, Value: v.Free}
	g.Metrics[KeyCPUutilization1] = Metrics{ID: CPUutilization1, MType: entities.Gauge, Value: int64(cpuNum)}

	g.mu.Unlock()
}

// GetMetricToSend возвращает массив метрик, готовых к отправке.
func (g *Gatherer) GetMetricToSend() []Metrics {
	g.mu.RLock()
	defer g.mu.RUnlock()
	metrics := make([]Metrics, 0, len(g.Metrics))
	for _, metric := range g.Metrics {
		metrics = append(metrics, metric)
	}

	return metrics
}

// StartMetricsGatherer запускает сбор метрик в фоновом режиме.
func (g *Gatherer) StartMetricsGatherer() {
	ticker := time.NewTicker(g.pollInterval)
	defer ticker.Stop()
	var delta int64

	for range ticker.C {
		go g.GatherRuntimeMetrics(&delta)
		go g.GatherMemoryMetrics()
	}
}
