package zeus_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func (z *ZeusClient) UpdateTopologyKnsStatus(ctx context.Context, status topology_deployment_status.Status) (zeus_req_types.TopologyDeployRequest, error) {
	z.PrintReqJson(status)
	respStatus := zeus_req_types.TopologyDeployRequest{}
	resp, err := z.R().
		SetResult(&respStatus).
		SetBody(status.TopologyDeployRequest).
		Post(zeus_endpoints.InternalDeployKnsCreateOrUpdatePath)

	if err != nil {
		log.Err(err).Msg("ZeusClient: UpdateTopologyKnsStatus")
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return respStatus, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		log.Err(err).Interface("status", status).Msg("ZeusClient: RemoveTopologyKnsStatus")
		return respStatus, err
	}
	z.PrintRespJson(resp.Body())
	return respStatus, err
}

func (z *ZeusClient) RemoveTopologyKnsStatus(ctx context.Context, status topology_deployment_status.Status) (zeus_req_types.TopologyDeployRequest, error) {
	z.PrintReqJson(status)
	respStatus := zeus_req_types.TopologyDeployRequest{}
	resp, err := z.R().
		SetResult(&respStatus).
		SetBody(status.TopologyDeployRequest).
		Post(zeus_endpoints.InternalDeployKnsDestroyPath)

	if err != nil {
		log.Err(err).Interface("status", status).Msg("ZeusClient: RemoveTopologyKnsStatus")
		return respStatus, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		log.Err(err).Interface("status", status).Msg("ZeusClient: RemoveTopologyKnsStatus")
		return respStatus, err
	}
	z.PrintRespJson(resp.Body())
	return respStatus, err
}
