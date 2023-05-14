package types

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric struct {
	MetricType   string
	MetricName   string
	GaugeValue   float64
	CounterValue int64
}
