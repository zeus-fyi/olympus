package athena_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	athena_endpoints "github.com/zeus-fyi/olympus/pkg/athena/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	zeus_pods_resp "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_resp_types/pods"
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

var AthenaPorts = []string{"9003:9003"}

const Athena = "athena"

func (a *AthenaClient) UploadViaPortForward(ctx context.Context, routeHeader zeus_req_types.TopologyDeployRequest, br poseidon.BucketRequest) (zeus_pods_resp.ClientResp, error) {
	cliReq := zeus_pods_reqs.ClientRequest{
		MethodHTTP: "POST",
		Endpoint:   athena_endpoints.InternalUploadV1Path,
		Ports:      AthenaPorts,
		Payload:    br,
	}
	filter := string_utils.FilterOpts{Contains: br.ClientName}
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: routeHeader,
		Action:                zeus_pods_reqs.PortForwardToAllMatchingPods,
		ContainerName:         Athena,
		ClientReq:             &cliReq,
		FilterOpts:            &filter,
	}
	resp, err := a.ZeusClient.PortForwardReqToPods(ctx, par)
	return resp, err
}
