package zeus_req_types

import "github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"

type TopologyRequest struct {
	TopologyID int `json:"topologyID"`
}

type TopologyCreateRequest struct {
	TopologyName     string `json:"topologyName"`
	ChartName        string `json:"chartName"`
	ChartDescription string `json:"chartDescription,omitempty"`
	Version          string `json:"version"`

	SkeletonBaseID int `json:"skeletonBaseID,omitempty"`
}

type TopologyDeployRequest struct {
	TopologyID int `json:"topologyID"`
	zeus_common_types.CloudCtxNs
}

type TopologyCloudCtxNsQueryRequest struct {
	zeus_common_types.CloudCtxNs
}

type TopologyCreateClusterRequest struct {
	ClusterName string   `json:"name"`
	Bases       []string `json:"bases,omitempty"`
}
