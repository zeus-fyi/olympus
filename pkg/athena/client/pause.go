package athena_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	athena_routines "github.com/zeus-fyi/olympus/athena/api/v1/common/routines"
	athena_endpoints "github.com/zeus-fyi/olympus/pkg/athena/client/endpoints"
)

func (a *AthenaClient) Pause(ctx context.Context, rr athena_routines.RoutineRequest) error {
	a.PrintReqJson(rr)
	resp, err := a.R().
		SetBody(rr).
		Post(athena_endpoints.InternalPauseV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("AthenaClient: Pause")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return err
	}
	a.PrintRespJson(resp.Body())
	return err
}
