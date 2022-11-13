package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

func (z *ZeusClient) DestroyDeploy(ctx context.Context, tar destroy_deploy_request.TopologyDestroyDeployRequest) (topology_deployment_status.Status, error) {
	z.PrintReqJson(tar)

	respJson := topology_deployment_status.Status{}
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(tar).
		Post(zeus_endpoints.DestroyDeployInfraV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: DestroyDeploy")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}
