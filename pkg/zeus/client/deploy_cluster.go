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

func (z *ZeusClient) DeployCluster(ctx context.Context, tar zeus_req_types.ClusterTopologyDeployRequest) (zeus_resp_types.ClusterStatus, error) {
	z.PrintReqJson(tar)
	respJson := zeus_resp_types.ClusterStatus{}
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(tar).
		Post(zeus_endpoints.DeployClusterTopologyV1Path)

	if err != nil || resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Msg("ZeusClient: DeployCluster")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}
