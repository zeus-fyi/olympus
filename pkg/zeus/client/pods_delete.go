package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
)

func (z *ZeusClient) DeletePods(ctx context.Context, par zeus_pods_reqs.PodActionRequest) ([]byte, error) {
	par.Action = zeus_pods_reqs.DeleteAllPods
	resp, err := z.R().
		SetBody(par).
		Post(zeus_endpoints.PodsActionV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: DeletePods")
		return resp.Body(), err
	}
	z.PrintRespJson(resp.Body())
	return resp.Body(), err
}
