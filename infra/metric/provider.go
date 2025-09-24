package metric

type MetricType string

const (
	MetricTypeCounter     MetricType = "counter"
	MetricTypeRateCounter MetricType = "rate_counter"
	MetricTypeStore       MetricType = "store"
	MetricTypeTimer       MetricType = "timer"
	MetricTypeHistogram   MetricType = "histogram"
)

//go:generate mockgen -destination ./mocks/provider.go -package=mocks . IMeter
type IMeter interface {
	NewMetric(name string, types []MetricType, tagNames []string) (IMetric, error)
}

var provider IMeter = noopMeter{}

// GetMeter Get the metric provider. Must call InitMeter first.
func GetMeter() IMeter {
	return provider
}

// InitMeter Init the metric provider. Must call before GetMeter.
func InitMeter(p IMeter) {
	provider = p
}
