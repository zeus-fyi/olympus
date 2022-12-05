package apollo_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	apollo_endpoints "github.com/zeus-fyi/olympus/pkg/apollo/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/apollo/client/req_types"
	"github.com/zeus-fyi/olympus/pkg/apollo/client/resp_types"
)

func (a *Apollo) ValidatorBalanceSums(ctx context.Context, rr req_types.ValidatorBalancesRequest) (resp_types.ValidatorBalancesSum, error) {
	a.PrintReqJson(rr)

	respJson := resp_types.ValidatorBalancesSum{}
	resp, err := a.R().
		SetBody(rr).
		SetResult(&respJson).
		Post(apollo_endpoints.EthereumValidatorsBalanceSumsV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("Apollo: ValidatorBalanceSums")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return respJson, err
	}

	a.PrintRespJson(resp.Body())
	return respJson, err
}
