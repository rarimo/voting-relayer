package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const configName = "autorelayer"

type AutorelayerConfiger interface {
	AutorelayConfig() *AutorelayerConfig
}

type AutorelayerConfig struct {
	CatchupDisabled bool `fig:"catchup_disabled"`
}

type autorelayer struct {
	once   comfig.Once
	getter kv.Getter
}

func NewAutorelayerConfiger(getter kv.Getter) AutorelayerConfiger {
	return &autorelayer{
		getter: getter,
	}
}

func (c *autorelayer) AutorelayConfig() *AutorelayerConfig {
	return c.once.Do(func() interface{} {
		var cfg AutorelayerConfig

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, configName)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out "+configName))
		}

		return &cfg
	}).(*AutorelayerConfig)
}
