package pods

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func PodsDeleteRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Debug().Msg("PodsDeleteRequest")

	err := topology_worker.Worker.ExecuteDeletePodWorkflow(ctx, request.CloudCtxNs, request.PodName, request.Delay)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsDeleteRequest: ExecuteDeletePodWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func PodsDeleteAllRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodsDeleteAllRequest")
	err := zeus.K8Util.DeleteAllPodsLike(ctx, request.CloudCtxNs, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsDeleteAllRequest: DeleteAllPodsLike")
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}
