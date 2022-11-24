package zeus_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_resp_types"
)

func (z *ZeusClient) Deploy(ctx context.Context, tar zeus_req_types.TopologyDeployRequest) (zeus_resp_types.TopologyDeployStatus, error) {
	z.PrintReqJson(tar)
	respJson := zeus_resp_types.TopologyDeployStatus{}
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(tar).
		Post(zeus_endpoints.DeployTopologyV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: Deploy")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}
