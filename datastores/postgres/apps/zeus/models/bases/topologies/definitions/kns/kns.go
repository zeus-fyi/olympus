package kns

import (
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyKubeCtxNs struct {
	TopologyID                   int `json:"topologyID"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
}

func NewKns() TopologyKubeCtxNs {
	k := TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "",
			Region:        "",
			Context:       "",
			Namespace:     "",
			Env:           "",
		},
	}
	return k
}

func NewTopologyCloudCtxNs(topID int, ccns zeus_common_types.CloudCtxNs) TopologyKubeCtxNs {
	return TopologyKubeCtxNs{
		TopologyID: topID,
		CloudCtxNs: ccns,
	}
}
