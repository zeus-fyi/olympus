package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
)

func (z *ZeusClient) UpdateTopologyStatus(ctx context.Context, status topology_deployment_status.Status) (topology_deployment_status.DeployStatus, error) {
	z.PrintReqJson(status)

	respStatus := topology_deployment_status.DeployStatus{}
	resp, err := z.R().
		SetResult(&respStatus).
		SetBody(status).
		Post(zeus_endpoints.InternalDeployStatusUpdatePath)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: UpdateTopologyStatus")
		return respStatus, err
	}
	z.PrintRespJson(resp.Body())
	return respStatus, err
}
