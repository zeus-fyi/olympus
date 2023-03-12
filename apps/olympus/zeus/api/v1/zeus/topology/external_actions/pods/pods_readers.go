package pods

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	v1 "k8s.io/api/core/v1"
)

func PodsDescribeRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	pods, err := zeus.K8Util.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsDescribeRequest")
		return err
	}
	return c.JSON(http.StatusOK, pods)
}

func PodLogsActionRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodLogsActionRequest")
	pods, err := zeus.K8Util.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, nil, request.FilterOpts)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodLogsActionRequest: GetPodsUsingCtxNs")
		return err
	}
	p := v1.Pod{}
	for _, pod := range pods.Items {
		if string_utils.FilterStringWithOpts(pod.GetName(), request.FilterOpts) {
			p = pod
		}
	}
	logs, err := zeus.K8Util.GetPodLogs(ctx, p.GetName(), request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodLogsActionRequest: GetPodLogs")
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}

func PodsAuditRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	pods, err := zeus.K8Util.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PodsAuditRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	parsedResp := parseResp(pods)
	return c.JSON(http.StatusOK, parsedResp)
}
