package gatherer

import (
	"math/rand"
	"runtime"
	"time"

	"metric-alert/internal/server/entities"
)

type Gatherer struct {
	Metrics      map[string]Metrics
	pollInterval time.Duration
}

func NewGatherer(pollInterval int) *Gatherer {
	return &Gatherer{
		Metrics:      make(map[string]Metrics, 0),
		pollInterval: time.Duration(pollInterval) * time.Second,
	}
}

func (a *Gatherer) GatherMetrics() {
	var m runtime.MemStats
	var delta int64
	for {
		runtime.ReadMemStats(&m)
		a.Metrics[Alloc] = Metrics{ID: Alloc, MType: entities.Gauge, Value: m.Alloc}
		a.Metrics[BuckHashSys] = Metrics{ID: BuckHashSys, MType: entities.Gauge, Value: m.BuckHashSys}
		a.Metrics[Frees] = Metrics{ID: Frees, MType: entities.Gauge, Value: m.Frees}
		a.Metrics[GCSys] = Metrics{ID: GCSys, MType: entities.Gauge, Value: m.GCSys}
		a.Metrics[HeapAlloc] = Metrics{ID: HeapAlloc, MType: entities.Gauge, Value: m.HeapAlloc}
		a.Metrics[HeapIdle] = Metrics{ID: HeapIdle, MType: entities.Gauge, Value: m.HeapIdle}
		a.Metrics[HeapInuse] = Metrics{ID: HeapInuse, MType: entities.Gauge, Value: m.HeapInuse}
		a.Metrics[HeapObjects] = Metrics{ID: HeapObjects, MType: entities.Gauge, Value: m.HeapObjects}
		a.Metrics[HeapReleased] = Metrics{ID: HeapReleased, MType: entities.Gauge, Value: m.HeapReleased}
		a.Metrics[HeapSys] = Metrics{ID: HeapSys, MType: entities.Gauge, Value: m.HeapSys}
		a.Metrics[LastGC] = Metrics{ID: LastGC, MType: entities.Gauge, Value: m.LastGC}
		a.Metrics[Lookups] = Metrics{ID: Lookups, MType: entities.Gauge, Value: m.Lookups}
		a.Metrics[MCacheInuse] = Metrics{ID: MCacheInuse, MType: entities.Gauge, Value: m.MCacheInuse}
		a.Metrics[MCacheSys] = Metrics{ID: MCacheSys, MType: entities.Gauge, Value: m.MCacheSys}
		a.Metrics[MSpanInuse] = Metrics{ID: MSpanInuse, MType: entities.Gauge, Value: m.MSpanInuse}
		a.Metrics[MSpanSys] = Metrics{ID: MSpanSys, MType: entities.Gauge, Value: m.MSpanSys}
		a.Metrics[Mallocs] = Metrics{ID: Mallocs, MType: entities.Gauge, Value: m.Mallocs}
		a.Metrics[NextGC] = Metrics{ID: NextGC, MType: entities.Gauge, Value: m.NextGC}
		a.Metrics[NumForcedGC] = Metrics{ID: NumForcedGC, MType: entities.Gauge, Value: m.NumForcedGC}
		a.Metrics[NumGC] = Metrics{ID: NumGC, MType: entities.Gauge, Value: m.NumGC}
		a.Metrics[OtherSys] = Metrics{ID: OtherSys, MType: entities.Gauge, Value: m.OtherSys}
		a.Metrics[PauseTotalNs] = Metrics{ID: PauseTotalNs, MType: entities.Gauge, Value: m.PauseTotalNs}
		a.Metrics[StackInuse] = Metrics{ID: StackInuse, MType: entities.Gauge, Value: m.StackInuse}
		a.Metrics[StackSys] = Metrics{ID: StackSys, MType: entities.Gauge, Value: m.StackSys}
		a.Metrics[Sys] = Metrics{ID: Sys, MType: entities.Gauge, Value: m.Sys}
		a.Metrics[TotalAlloc] = Metrics{ID: TotalAlloc, MType: entities.Gauge, Value: m.TotalAlloc}
		a.Metrics[GCCPUFraction] = Metrics{ID: GCCPUFraction, MType: entities.Gauge, Value: m.GCCPUFraction}
		a.Metrics[RandomValue] = Metrics{ID: RandomValue, MType: entities.Gauge, Value: rand.Uint32()}
		a.Metrics[PollCount] = Metrics{ID: PollCount, MType: entities.Counter, Delta: delta}

		delta++
		time.Sleep(a.pollInterval)
	}
}
