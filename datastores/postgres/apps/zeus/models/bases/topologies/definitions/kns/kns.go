package kns

import (
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type TopologyKubeCtxNs struct {
	TopologyID int `db:"topology_id" json:"topologyID"`
	zeus_core.CloudCtxNs
}

func NewKns() TopologyKubeCtxNs {
	k := TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: zeus_core.NewCloudCtxNs(),
	}
	return k
}

func NewTopologyCloudCtxNs(topID int, ccns zeus_core.CloudCtxNs) TopologyKubeCtxNs {
	return TopologyKubeCtxNs{
		TopologyID: topID,
		CloudCtxNs: ccns,
	}
}
