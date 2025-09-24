package metric

type noopMeter struct{}

func (n noopMeter) NewMetric(name string, types []MetricType, tagNames []string) (IMetric, error) {
	return noopMetric{}, nil
}

type noopMetric struct{}

func (n noopMetric) Emit(tags []Tag, values ...*Value) {}
