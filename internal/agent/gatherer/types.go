package gatherer

const (
	DefaultMetricCount = 32

	KeyAlloc = iota
	KeyBuckHashSys
	KeyFrees
	KeyGCCPUFraction
	KeyGCSys
	KeyHeapAlloc
	KeyHeapIdle
	KeyHeapInuse
	KeyHeapObjects
	KeyHeapReleased
	KeyHeapSys
	KeyLastGC
	KeyLookups
	KeyMCacheInuse
	KeyMCacheSys
	KeyMSpanInuse
	KeyMSpanSys
	KeyMallocs
	KeyNextGC
	KeyNumForcedGC
	KeyNumGC
	KeyOtherSys
	KeyPauseTotalNs
	KeyStackInuse
	KeyStackSys
	KeySys
	KeyTotalAlloc
	KeyRandomValue
	KeyPollCount
	KeyTotalMemory
	KeyFreeMemory
	KeyCPUutilization1

	Alloc           = "Alloc"
	BuckHashSys     = "BuckHashSys"
	Frees           = "Frees"
	GCCPUFraction   = "GCCPUFraction"
	GCSys           = "GCSys"
	HeapAlloc       = "HeapAlloc"
	HeapIdle        = "HeapIdle"
	HeapInuse       = "HeapInuse"
	HeapObjects     = "HeapObjects"
	HeapReleased    = "HeapReleased"
	HeapSys         = "HeapSys"
	LastGC          = "LastGC"
	Lookups         = "Lookups"
	MCacheInuse     = "MCacheInuse"
	MCacheSys       = "MCacheSys"
	MSpanInuse      = "MSpanInuse"
	MSpanSys        = "MSpanSys"
	Mallocs         = "Mallocs"
	NextGC          = "NextGC"
	NumForcedGC     = "NumForcedGC"
	NumGC           = "NumGC"
	OtherSys        = "OtherSys"
	PauseTotalNs    = "PauseTotalNs"
	StackInuse      = "StackInuse"
	StackSys        = "StackSys"
	Sys             = "Sys"
	TotalAlloc      = "TotalAlloc"
	RandomValue     = "RandomValue"
	PollCount       = "PollCount"
	TotalMemory     = "TotalMemory"
	FreeMemory      = "FreeMemory"
	CPUutilization1 = "CPUutilization1"
)

type Metrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Value any    `json:"value,omitempty"`
	Delta int64  `json:"delta,omitempty"`
}
