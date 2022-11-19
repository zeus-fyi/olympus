package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	zeus_pods_resp "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_resp_types/pods"
)

// PortForwardReqToPods will use this filter by default if you specify a pod name without your own custom filter override
// filter request.FilterOpts.StartsWith = request.PodName
func (z *ZeusClient) PortForwardReqToPods(ctx context.Context, par zeus_pods_reqs.PodActionRequest) (zeus_pods_resp.ClientResp, error) {
	par.Action = zeus_pods_reqs.PortForwardToAllMatchingPods

	clientResponses := zeus_pods_resp.ClientResp{}
	resp, err := z.R().
		SetResult(&clientResponses).
		SetBody(par).
		Post(zeus_endpoints.PodsActionV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: PortForwardReqToPods")
		return clientResponses, err
	}
	z.PrintRespJson(resp.Body())
	return clientResponses, err
}
