package pods

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func HandlePodActionRequest(c echo.Context) error {
	request := new(PodActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)

	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
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
		bytesResp, perr := podsPortForwardRequest(request)
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
