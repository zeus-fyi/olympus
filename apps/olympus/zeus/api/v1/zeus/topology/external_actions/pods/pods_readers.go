package pods

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
	v1 "k8s.io/api/core/v1"
)

func PodsDescribeRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	ctx := context.Background()
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}
	pods, err := k.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Err(err).Msg("PodsDescribeRequest")
		return err
	}
	return c.JSON(http.StatusOK, pods)
}

func PodLogsActionRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	ctx := context.Background()
	log.Debug().Msg("PodLogsActionRequest")
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}

	pods, err := k.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, nil, request.FilterOpts)
	if err != nil {
		log.Err(err).Msg("PodLogsActionRequest: GetPodsUsingCtxNs")
		return err
	}
	p := v1.Pod{}
	for _, pod := range pods.Items {
		if strings_filter.FilterStringWithOpts(pod.GetName(), request.FilterOpts) {
			p = pod
		}
	}
	if request.LogOpts == nil {
		request.LogOpts = &v1.PodLogOptions{
			Container: request.ContainerName,
		}
	}
	if request.ContainerName == "" {
		request.ContainerName = p.Spec.Containers[len(p.Spec.Containers)-1].Name
		request.LogOpts.Container = request.ContainerName
	}
	logs, err := k.GetPodLogs(ctx, p.GetName(), request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Err(err).Msg("PodLogsActionRequest: GetPodLogs")
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}

func PodsAuditRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	ctx := context.Background()
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}
	pods, err := k.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, request.LogOpts, request.FilterOpts)
	if err != nil {
		log.Err(err).Msg("PodsAuditRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	parsedResp := parseResp(pods)
	return c.JSON(http.StatusOK, parsedResp)
}
