package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/chaos-io/chaos/core/metrics"
)

var _ metrics.Histogram = (*Histogram)(nil)

type Histogram struct {
	hm prometheus.Observer
}

func (h Histogram) RecordValue(value float64) {
	h.hm.Observe(value)
}

func (h Histogram) RecordDuration(value time.Duration) {
	h.hm.Observe(value.Seconds())
}
