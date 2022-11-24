package athena_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	athena_endpoints "github.com/zeus-fyi/olympus/pkg/athena/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
)

func (a *AthenaClient) Upload(ctx context.Context, br poseidon.BucketRequest) error {
	a.PrintReqJson(br)
	resp, err := a.R().
		SetBody(br).
		Post(athena_endpoints.InternalUploadV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("AthenaClient: Upload")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return err
	}
	a.PrintRespJson(resp.Body())
	return err
}
