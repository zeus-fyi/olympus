package coreK8s

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func HandlePodActionRequest(c echo.Context) error {
	request := new(PodActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "logs" {
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
		bytesResp, err := podsPortForwardRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, string(bytesResp))
		}
		return c.JSON(http.StatusOK, string(bytesResp))
	}
	if request.Action == "port-forward-all" {
		return podsPortForwardRequestToAllPods(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
