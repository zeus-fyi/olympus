package kns

import (
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type TopologyKubeCtxNs struct {
	TopologyID int `db:"topology_id" json:"topologyID"`
	zeus_common_types.CloudCtxNs
}

func NewKns() TopologyKubeCtxNs {
	k := TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: zeus_common_types.NewCloudCtxNs(),
	}
	return k
}

func NewTopologyCloudCtxNs(topID int, ccns zeus_common_types.CloudCtxNs) TopologyKubeCtxNs {
	return TopologyKubeCtxNs{
		TopologyID: topID,
		CloudCtxNs: ccns,
	}
}
