package entities

import "fmt"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m Metrics) ByteValue() []byte {
	if m.MType == Gauge {
		return []byte(fmt.Sprintf("%v", *m.Value))
	} else {
		return []byte(fmt.Sprintf("%v", *m.Delta))
	}
}
