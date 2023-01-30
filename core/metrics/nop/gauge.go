package nop

import "github.com/chaos-io/chaos/core/metrics"

var _ metrics.Gauge = (*Gauge)(nil)

type Gauge struct{}

func (Gauge) Set(_ float64) {}

func (Gauge) Add(_ float64) {}

var _ metrics.GaugeVec = (*GaugeVec)(nil)

type GaugeVec struct{}

func (t GaugeVec) With(_ map[string]string) metrics.Gauge {
	return Gauge{}
}
