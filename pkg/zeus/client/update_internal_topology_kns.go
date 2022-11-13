package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
)

func (z *ZeusClient) UpdateTopologyKnsStatus(ctx context.Context, status topology_deployment_status.Status) (kns.TopologyKubeCtxNs, error) {
	z.PrintReqJson(status)
	respStatus := kns.TopologyKubeCtxNs{}
	resp, err := z.R().
		SetResult(&respStatus).
		SetBody(status.TopologyKubeCtxNs).
		Post(zeus_endpoints.InternalDeployKnsStatusUpdatePath)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: UpdateTopologyKnsStatus")
		return respStatus, err
	}
	z.PrintRespJson(resp.Body())
	return respStatus, err
}
