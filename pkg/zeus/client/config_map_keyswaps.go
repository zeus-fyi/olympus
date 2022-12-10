package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	zeus_configmap_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/config_maps"
)

func (z *ZeusClient) SwapConfigMapKeys(ctx context.Context, par zeus_configmap_reqs.ConfigMapActionRequest) ([]byte, error) {
	par.Action = zeus_configmap_reqs.KeySwapAction
	resp, err := z.R().
		SetBody(par).
		Post(zeus_endpoints.ConfigMapsActionV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: SwapConfigMapKeys")
		return resp.Body(), err
	}
	z.PrintRespJson(resp.Body())
	return resp.Body(), err
}
