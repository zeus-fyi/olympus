package zeus_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
)

func (z *ZeusClient) CopySecretsFromToNamespace(ctx context.Context, secretsCopyReq internal_reqs.InternalSecretsCopyFromTo) error {
	z.PrintReqJson(secretsCopyReq)

	resp, err := z.R().
		SetBody(secretsCopyReq).
		Post(zeus_endpoints.InternalSecretsCopyFromTo)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: CopySecretsFromToNamespace")
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return err
	}
	z.PrintRespJson(resp.Body())
	return err
}
