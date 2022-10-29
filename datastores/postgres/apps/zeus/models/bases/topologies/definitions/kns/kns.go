package kns

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type Kns struct {
	autogen_bases.TopologiesKns
}

func NewKns() Kns {
	k := Kns{autogen_bases.TopologiesKns{
		TopologyID: 0,
		Context:    "",
		Namespace:  "",
		Env:        "",
	}}
	return k
}
