package infra

import (
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	delete_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/delete"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/update"
)

type TopologyInfraRequest struct {
	Action string

	create_infra.TopologyActionCreateRequest
	read_infra.TopologyActionReadRequest
	update_infra.TopologyActionUpdateRequest
	delete_infra.TopologyActionDeleteRequest
}
