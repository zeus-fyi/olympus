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

	Tag              string `json:"tag,omitempty"`
	SkeletonBaseName string `json:"skeletonBaseName,omitempty"`
}

type TopologyDeployRequest struct {
	TopologyID int `json:"topologyID"`
	zeus_common_types.CloudCtxNs
}

type TopologyCloudCtxNsQueryRequest struct {
	zeus_common_types.CloudCtxNs
}

type TopologyCreateOrAddBasesToClassesRequest struct {
	ClassName      string   `json:"className"`
	ClassBaseNames []string `json:"classBaseNames,omitempty"`
}
