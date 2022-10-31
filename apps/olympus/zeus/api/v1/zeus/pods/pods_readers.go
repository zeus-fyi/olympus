package pods

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
	v1 "k8s.io/api/core/v1"
)

func PodsDescribeRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	pods, err := core.K8Util.GetPodsUsingCtxNs(ctx, request.Kns, request.LogOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pods)
}

func PodLogsActionRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodLogsActionRequest")
	pods, err := core.K8Util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
	if err != nil {
		return err
	}

	p := v1.Pod{}
	for _, pod := range pods.Items {
		if string_utils.FilterStringWithOpts(pod.GetName(), request.FilterOpts) {
			p = pod
		}
	}
	logs, err := core.K8Util.GetPodLogs(ctx, p.GetName(), request.Kns.Namespace, request.LogOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}

func PodsAuditRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()

	pods, err := core.K8Util.GetPodsUsingCtxNs(ctx, request.Kns, request.LogOpts, request.FilterOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	parsedResp := parseResp(pods)
	return c.JSON(http.StatusOK, parsedResp)
}
