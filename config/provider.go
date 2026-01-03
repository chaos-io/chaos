package config

//go:generate mockgen -destination=mocks/config_provider.go -package=mocks . IConfigProvider
type IConfigProvider interface {
	Get(key string) (any, error)
	Scan(v any) error
	ScanFrom(v any, key string, alternatives ...string) error
}

//go:generate mockgen -destination=mocks/config_provider_factory.go -package=mocks . IConfigProviderFactory
type IConfigProviderFactory interface {
	NewConfigProvider() IConfigProvider
}

type defaultConfigProvider struct{}

func NewDefaultConfigProvider() IConfigProvider {
	return &defaultConfigProvider{}
}

func (p *defaultConfigProvider) Get(key string) (any, error) {
	return Get(key)
}

func (p *defaultConfigProvider) Scan(v any) error {
	return Scan(v)
}

func (p *defaultConfigProvider) ScanFrom(v any, key string, alternatives ...string) error {
	return ScanFrom(v, key, alternatives...)
}
