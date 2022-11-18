package pods

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

func HandlePodActionRequest(c echo.Context) error {
	request := new(PodActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)

	tempKns := request.K8sRequest.Kns

	// TODO refactor
	knsDeploy := kns.NewKns()
	knsDeploy.TopologiesKns = autogen_bases.TopologiesKns{
		CloudProvider: tempKns.CloudProvider,
		Region:        tempKns.Region,
		Context:       tempKns.Context,
		Namespace:     tempKns.Namespace,
		Env:           tempKns.Env,
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, knsDeploy)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
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
		bytesResp, perr := podsPortForwardRequest(request)
		if perr != nil {
			return c.JSON(http.StatusBadRequest, string(bytesResp))
		}
		return c.JSON(http.StatusOK, string(bytesResp))
	}
	if request.Action == "port-forward-all" {
		return podsPortForwardRequestToAllPods(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
