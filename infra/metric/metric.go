package metric

//go:generate mockgen -destination=mocks/metric.go -package=mocks . IMetric
type IMetric interface {
	Emit(tags []Tag, values ...*Value)
}

type Tag struct {
	Name  string
	Value string
}

type Value struct {
	suffix string
	mType  MetricType
	value  *int64
	valueF *float64
}

func (v *Value) GetSuffix() string { return v.suffix }
func (v *Value) GetType() MetricType { return v.mType }
func (v *Value) GetValue() *int64 { return v.value }
func (v *Value) GetValueF() *float64 { return v.valueF }

type ValueOption func(*Value)

func WithSuffix(suffix string) ValueOption {
	return func(v *Value) {
		v.suffix = suffix
	}
}

func Counter(n int64, opts ...ValueOption) *Value { return buildValue(n, MetricTypeCounter, opts...) }
func RateCounter(n int64, opts ...ValueOption) *Value {return buildValue(n, MetricTypeRateCounter, opts...)}
func Store(n int64, opts ...ValueOption) *Value { return buildValue(n, MetricTypeStore, opts...) }
func Timer(n int64, opts ...ValueOption) *Value { return buildValue(n, MetricTypeTimer, opts...) }
func Histogram(n int64, opts ...ValueOption) *Value { return buildValue(n, MetricTypeHistogram, opts...) }

func CounterF(n float64, opts ...ValueOption) *Value { return buildValueF(n, MetricTypeCounter, opts...) }
func RateCounterF(n float64, opts ...ValueOption) *Value { return buildValueF(n, MetricTypeRateCounter, opts...) }
func StoreF(n float64, opts ...ValueOption) *Value { return buildValueF(n, MetricTypeStore, opts...) }
func TimerF(n float64, opts ...ValueOption) *Value { return buildValueF(n, MetricTypeTimer, opts...) }
func HistogramF(n float64, opts ...ValueOption) *Value { return buildValueF(n, MetricTypeHistogram, opts...) }

func buildValue(n int64, mType MetricType, opts ...ValueOption) *Value {
	v := &Value{
		mType: mType,
		value: &n,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func buildValueF(n float64, mType MetricType, opts ...ValueOption) *Value {
	value := &Value{
		mType:  mType,
		valueF: &n,
	}
	for _, opt := range opts {
		opt(value)
	}
	return value
}
