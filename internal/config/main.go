package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	pgdb.Databaser

	Tenderminter
	Cosmoser
	EVMer

	RelayerConfiger
	AutorelayerConfiger
}

type config struct {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
	pgdb.Databaser

	Tenderminter
	Cosmoser
	EVMer

	RelayerConfiger
	AutorelayerConfiger
}

func New(getter kv.Getter) Config {
	return &config{
		getter:              getter,
		Copuser:             copus.NewCopuser(getter),
		Listenerer:          comfig.NewListenerer(getter),
		Logger:              comfig.NewLogger(getter, comfig.LoggerOpts{}),
		RelayerConfiger:     NewRelayerConfiger(getter),
		Tenderminter:        NewTenderminter(getter),
		Cosmoser:            NewCosmoser(getter),
		EVMer:               NewEVMer(getter),
		Databaser:           pgdb.NewDatabaser(getter),
		AutorelayerConfiger: NewAutorelayerConfiger(getter),
	}
}
