package pods

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func HandlePodActionRequest(c echo.Context) error {
	request := c.Get("PodActionRequest").(*PodActionRequest)
	if request.Action == "logs" {
		if request.FilterOpts == nil {
			request.FilterOpts = &string_utils.FilterOpts{}
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
			request.FilterOpts = &string_utils.FilterOpts{}
			request.FilterOpts.StartsWith = request.PodName
		}
		return podsPortForwardRequestToAllPods(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
