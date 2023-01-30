package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/chaos-io/chaos/core/metrics"
)

var _ metrics.Gauge = (*Gauge)(nil)

// Gauge tracks single float64 value.
type Gauge struct {
	gg prometheus.Gauge
}

func (g Gauge) Set(value float64) {
	g.gg.Set(value)
}

func (g Gauge) Add(value float64) {
	g.gg.Add(value)
}
