package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

func (z *ZeusClient) UpdateTopologyStatus(ctx context.Context, status topology_deployment_status.Status) (topology_deployment_status.Status, error) {
	z.PrintReqJson(status)

	respJson := topology_deployment_status.Status{}
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(status).
		Post(InternalDeployStatusUpdatePath)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: UpdateTopologyStatus")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}
