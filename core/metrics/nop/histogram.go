package nop

import (
	"time"

	"github.com/chaos-io/chaos/core/metrics"
)

var (
	_ metrics.Histogram = (*Histogram)(nil)
	_ metrics.Timer     = (*Histogram)(nil)
)

type Histogram struct{}

func (Histogram) RecordValue(_ float64) {}

func (Histogram) RecordDuration(_ time.Duration) {}

var _ metrics.HistogramVec = (*HistogramVec)(nil)

type HistogramVec struct{}

func (t HistogramVec) With(_ map[string]string) metrics.Histogram {
	return Histogram{}
}

var _ metrics.TimerVec = (*DurationHistogramVec)(nil)

type DurationHistogramVec struct{}

func (t DurationHistogramVec) With(_ map[string]string) metrics.Timer {
	return Histogram{}
}
