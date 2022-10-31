package base_infra

import "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"

type TopologyInfraActionRequest struct {
	base.TopologyActionRequest

	TopologyCreateRequest
}

type TopologyCreateRequest struct {
	Name string
}
