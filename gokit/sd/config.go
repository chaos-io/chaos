package sd

import (
	"strings"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/logs"

	"github.com/chaos-io/chaos/gokit/retry"
	"github.com/chaos-io/chaos/gokit/sd/direct"
	"github.com/chaos-io/chaos/gokit/sd/etcdv3"
	"github.com/chaos-io/chaos/gokit/sd/nacos"
)

type Config struct {
	Mode   string                    `json:"mode" yaml:"mode" db:"mode"`
	Url    string                    `json:"url" yaml:"url"`
	Retry  *retry.Config             `json:"retry" yaml:"retry" db:"retry"`
	EtcdV3 *etcdv3.Config            `json:"etcd" yaml:"etcd"`
	Nacos  *nacos.Config             `json:"nacos" yaml:"nacos"`
	Direct map[string]*direct.Config `json:"direct" yaml:"direct" db:"direct"`
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}
	if err := config.ScanFrom(&cfg, "sd"); err != nil {
		logs.Errorw("failed to get the sd config from "+strings.Join(path, "."), "error", err)
		return nil
	}
	return cfg
}
