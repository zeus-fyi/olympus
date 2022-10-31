package pods

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func PodsDeleteRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodsDeleteRequest")
	err := core.K8Util.DeleteFirstPodLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func PodsDeleteAllRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodsDeleteAllRequest")
	err := core.K8Util.DeleteAllPodsLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}
