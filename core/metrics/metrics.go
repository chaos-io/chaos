// Package metrics provides interface collecting performance metrics.
//
//
package metrics

import (
	"context"
	"time"
)

// Gauge tracks single float64 value.
type Gauge interface {
	Set(value float64)
	Add(value float64)
}

// Counter tracks monotonically increasing value.
type Counter interface {
	// Inc increments counter by 1.
	Inc()

	// Add adds delta to the counter. Delta must be >=0.
	Add(delta int64)
}

// Histogram tracks distribution of value.
type Histogram interface {
	RecordValue(value float64)
}

// Timer measures durations.
type Timer interface {
	RecordDuration(value time.Duration)
}

// DurationBuckets defines buckets of the duration histogram.
type DurationBuckets interface {
	// Size returns number of buckets.
	Size() int

	// MapDuration returns index of the bucket.
	//
	// index is integer in range [0, Size()).
	MapDuration(d time.Duration) int

	// UpperBound of the last bucket is always +Inf.
	//
	// bucketIndex is integer in range [0, Size()-1).
	UpperBound(bucketIndex int) time.Duration
}

// Buckets defines intervals of the regular histogram.
type Buckets interface {
	// Size returns number of buckets.
	Size() int

	// MapValue returns index of the bucket.
	//
	// Index is integer in range [0, Size()).
	MapValue(v float64) int

	// UpperBound of the last bucket is always +Inf.
	//
	// bucketIndex is integer in range [0, Size()-1)
	UpperBound(bucketIndex int) float64
}

// GaugeVec stores multiple dynamically created gauges
type GaugeVec interface {
	With(map[string]string) Gauge
}

// CounterVec stores multiple dynamically created counters
type CounterVec interface {
	With(map[string]string) Counter
}

// TimerVec stores multiple dynamically created timers
type TimerVec interface {
	With(map[string]string) Timer
}

// HistogramVec stores multiple dynamically created histograms
type HistogramVec interface {
	With(map[string]string) Histogram
}

// Registry creates profiling metrics.
type Registry interface {
	// WithTags creates new sub-scope, where each metric has tags attached to it.
	WithTags(tags map[string]string) Registry
	// WithPrefix creates new sub-scope, where each metric has prefix added to it name.
	WithPrefix(prefix string) Registry

	ComposeName(parts ...string) string

	Counter(name string) Counter
	CounterVec(name string, labels []string) CounterVec
	FuncCounter(name string, function func() int64)

	Gauge(name string) Gauge
	GaugeVec(name string, labels []string) GaugeVec
	FuncGauge(name string, function func() float64)

	Timer(name string) Timer
	TimerVec(name string, labels []string) TimerVec

	Histogram(name string, buckets Buckets) Histogram
	HistogramVec(name string, buckets Buckets, labels []string) HistogramVec

	DurationHistogram(name string, buckets DurationBuckets) Timer
	DurationHistogramVec(name string, buckets DurationBuckets, labels []string) TimerVec
}

// CollectPolicy defines how registered gauge metrics are updated via collect func.
type CollectPolicy interface {
	RegisteredCounter(counterFunc func() int64) func() int64
	RegisteredGauge(gaugeFunc func() float64) func() float64
	AddCollect(collect func(ctx context.Context))
}
