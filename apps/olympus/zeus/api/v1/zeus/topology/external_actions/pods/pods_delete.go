package pods

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	pods_workflows "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/pods"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
)

func PodsDeleteRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	log.Debug().Msg("PodsDeleteRequest")

	err := pods_workflows.ExecuteDeletePodWorkflow(c, context.Background(), request.CloudCtxNs, request.PodName, request.Delay)
	if err != nil {
		log.Err(err).Msg("PodsDeleteRequest: ExecuteDeletePodWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func PodsDeleteAllRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	ctx := context.Background()
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}
	err := k.DeleteAllPodsLike(ctx, request.CloudCtxNs, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		log.Err(err).Msg("PodsDeleteAllRequest: DeleteAllPodsLike")
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}
