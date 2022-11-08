package kns

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type TopologyKubeCtxNs struct {
	autogen_bases.TopologiesKns
}

func NewKns() TopologyKubeCtxNs {
	k := TopologyKubeCtxNs{autogen_bases.TopologiesKns{
		TopologyID:    0,
		CloudProvider: "",
		Region:        "",
		Context:       "",
		Namespace:     "",
		Env:           "",
	}}
	return k
}
