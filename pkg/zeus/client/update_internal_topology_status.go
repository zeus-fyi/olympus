package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/deploy/workload_state"
)

func (z *ZeusClient) UpdateTopologyStatus(ctx context.Context, tar workload_state.InternalWorkloadStatusUpdateRequest) (topology_deployment_status.Status, error) {
	z.PrintReqJson(tar)

	respJson := topology_deployment_status.Status{}
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(tar).
		Post(InternalDeployStatusUpdatePath)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: UpdateTopologyStatus")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}
