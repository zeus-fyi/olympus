package snapshot_poseidon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	athena_endpoints "github.com/zeus-fyi/olympus/pkg/athena/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/external_actions/pods"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type InternalDeploymentActionRequest struct {
	ProtocolName  string
	ConfigMapName string
	Kns           kns.TopologyKubeCtxNs
}

func SnapshotHandler(c echo.Context) error {
	request := new(InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return request.SnapshotProcedure(c)
}

func (s *InternalDeploymentActionRequest) SnapshotProcedure(c echo.Context) error {
	request := new(InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	ctx := context.Background()
	_, err := zeus.K8Util.ConfigMapKeySwap(ctx, request.Kns.CloudCtxNs, request.ConfigMapName, "start.sh", "pause.sh", nil)
	if err != nil {
		log.Err(err).Msg("SnapshotProcedure")
		return c.JSON(http.StatusInternalServerError, err)
	}
	br := poseidon.BucketRequest{}
	br.CompressionType = "none"
	filter := string_utils.FilterOpts{DoesNotInclude: []string{""}}

	switch request.ProtocolName {
	case "geth":
		br = poseidon_buckets.GethMainnetBucket
		filter = string_utils.FilterOpts{DoesNotInclude: []string{"lighthouse"}}
	case "lighthouse":
		br = poseidon_buckets.LighthouseMainnetBucket
		filter = string_utils.FilterOpts{DoesNotInclude: []string{"geth"}}
	}
	err = zeus.K8Util.DeleteFirstPodLike(ctx, request.Kns.CloudCtxNs, fmt.Sprintf("zeus-%s-0", request.ProtocolName), nil, &filter)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsDeleteRequest: DeleteFirstPodLike")
		return err
	}

	cliReq := pods.ClientRequest{
		MethodHTTP: "POST",
		Endpoint:   athena_endpoints.InternalUploadV1Path,
		Ports:      []string{"9003:9003"},
		Payload:    br,
	}

	par := pods.PodActionRequest{
		Action:        zeus_pods_reqs.PortForwardToAllMatchingPods,
		PodName:       fmt.Sprintf("zeus-%s-0", request.ProtocolName),
		ContainerName: request.ProtocolName,
		ClientReq:     &cliReq,
		FilterOpts:    &filter,
	}

	_, err = pods.PodsPortForwardRequest(c, &par)
	if err != nil {
		log.Err(err).Msg("SnapshotProcedure")
		return c.JSON(http.StatusInternalServerError, err)
	}
	_, _ = zeus.K8Util.ConfigMapKeySwap(ctx, request.Kns.CloudCtxNs, request.ConfigMapName, "start.sh", "pause.sh", nil)
	err = zeus.K8Util.DeleteFirstPodLike(ctx, request.Kns.CloudCtxNs, fmt.Sprintf("zeus-%s-0", request.ProtocolName), nil, &filter)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsDeleteRequest: DeleteFirstPodLike")
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
