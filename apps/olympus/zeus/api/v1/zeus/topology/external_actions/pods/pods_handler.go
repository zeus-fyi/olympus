package pods

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
)

func HandlePodActionRequest(c echo.Context) error {
	request := c.Get("PodActionRequest").(*zeus_pods_reqs.PodActionRequest)
	if request.Action == "logs" {
		if request.FilterOpts == nil {
			request.FilterOpts = &strings_filter.FilterOpts{}
			request.FilterOpts.StartsWith = request.PodName
		}
		return PodLogsActionRequest(c, request)
	}
	if request.Action == "describe" {
		return PodsDescribeRequest(c, request)
	}
	if request.Action == "describe-audit" {
		return PodsAuditRequest(c, request)
	}
	if request.Action == "delete" {
		return PodsDeleteRequest(c, request)
	}
	if request.Action == "delete-all" {
		return PodsDeleteAllRequest(c, request)
	}
	if request.Action == "delete-all-delay" {
		time.Sleep(time.Second * 180)
		return PodsDeleteAllRequest(c, request)
	}
	if request.Action == "port-forward" {
		bytesResp, perr := PodsPortForwardRequest(c, request)
		if perr != nil {
			return c.JSON(http.StatusBadRequest, string(bytesResp))
		}
		return c.JSON(http.StatusOK, string(bytesResp))
	}
	if request.Action == "port-forward-all" {
		if request.FilterOpts == nil && len(request.PodName) > 0 {
			request.FilterOpts = &strings_filter.FilterOpts{}
			request.FilterOpts.StartsWith = request.PodName
		}
		return podsPortForwardRequestToAllPods(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
