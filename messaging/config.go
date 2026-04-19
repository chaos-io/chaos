package messaging

import (
	"strings"

	"github.com/chaos-io/chaos/config"
)

const (
	DriverNATS       = "nats"
	defaultConfigKey = "messaging"
)

type Config struct {
	Driver        string               `json:"driver"`
	Nats          NatsConfig           `json:"nats"`
	Subscriptions []SubscriptionConfig `json:"subscriptions"`
}

type NatsConfig struct {
	URL       string `json:"url" default:"nats://127.0.0.1:4222"`
	JetStream bool   `json:"jetStream"`
}

type SubscriptionConfig struct {
	Name    string `json:"name"`
	Topic   string `json:"topic"`
	Group   string `json:"group"`
	Service string `json:"service"`
	Method  string `json:"method"`
	Pull    bool   `json:"pull"`
	AutoAck bool   `json:"autoAck"`
	AckWait string `json:"ackWait"`
}

func New() (*Client, error) {
	cfg := &Config{}
	if err := config.ScanFrom(cfg, defaultConfigKey); err != nil {
		return nil, err
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *Config) (*Client, error) {
	queue, err := buildQueue(cfg)
	if err != nil {
		return nil, err
	}

	return NewClient(queue)
}

func (c *Config) normalized() (*Config, error) {
	if c == nil {
		return nil, ErrNilConfig
	}

	cfg := *c
	cfg.Driver = strings.ToLower(strings.TrimSpace(cfg.Driver))
	if cfg.Driver == "" {
		cfg.Driver = DriverNATS
	}

	return &cfg, nil
}
