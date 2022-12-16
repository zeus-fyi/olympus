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

	Tag               string `json:"tag,omitempty"`
	ClusterBaseName   string `json:"clusterBaseName,omitempty"`
	ComponentBaseName string `json:"componentBaseName,omitempty"`
	SkeletonBaseName  string `json:"skeletonBaseName,omitempty"`
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

type ClusterTopologyDeployRequest struct {
	ClusterName string   `json:"clusterName"`
	BaseOptions []string `json:"baseOptions"`
	zeus_common_types.CloudCtxNs
}

type ClusterTopology struct {
	ClusterName string              `json:"clusterName"`
	Topologies  []ClusterTopologies `json:"topologies"`
}

type ClusterTopologies struct {
	TopologyID       int    `json:"topologyID"`
	SkeletonBaseName string `json:"skeletonBaseName"`
	Tag              string `json:"tag"`
}
