package zeus_req_types

import "github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"

type TopologyRequest struct {
	TopologyID int `json:"topologyID"`
}

type TopologyCreateRequest struct {
	TopologyName     string `json:"topologyName"`
	ChartName        string `json:"chartName"`
	ChartDescription string `json:"chartDescription,omitempty"`
	Version          string `json:"version"`

	Tag               string `json:"tag,omitempty"`
	ClusterClassName  string `json:"clusterClassName,omitempty"`
	ComponentBaseName string `json:"componentBaseName,omitempty"`
	SkeletonBaseName  string `json:"skeletonBaseName,omitempty"`
}

type TopologyDeployRequest struct {
	TopologyID                   int `json:"topologyID"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`

	RequestChoreographySecretDeploy bool `json:"requestChoreographySecretDeploy"`
}

type TopologyCloudCtxNsQueryRequest struct {
	zeus_common_types.CloudCtxNs
}

type TopologyCreateClusterClassRequest struct {
	ClusterClassName string `json:"clusterClassName"`
}

type TopologyCreateOrAddComponentBasesToClassesRequest struct {
	ClusterClassName   string   `json:"clusterClassName,omitempty"`
	ComponentBaseNames []string `json:"componentBaseNames,omitempty"`
}

type TopologyCreateOrAddSkeletonBasesToClassesRequest struct {
	ClusterClassName  string   `json:"clusterClassName"`
	ComponentBaseName string   `json:"componentBaseName,omitempty"`
	SkeletonBaseNames []string `json:"skeletonBaseNames,omitempty"`
}

type ClusterTopologyDeployRequest struct {
	ClusterClassName             string   `json:"clusterClassName"`
	SkeletonBaseOptions          []string `json:"skeletonBaseOptions"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
}

type ClusterTopology struct {
	ClusterClassName string              `json:"clusterClassName"`
	Topologies       []ClusterTopologies `json:"topologies"`
}

type ClusterTopologies struct {
	TopologyID       int    `json:"topologyID"`
	SkeletonBaseName string `json:"skeletonBaseName"`
	Tag              string `json:"tag"`
}
