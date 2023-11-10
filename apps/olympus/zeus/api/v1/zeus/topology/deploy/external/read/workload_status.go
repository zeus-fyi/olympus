package deployment_status

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
)

type TopologyDeploymentStatusRequest struct {
	TopologyID int `json:"topologyID"`
}

func (t *TopologyDeploymentStatusRequest) ReadDeployedTopologyStatuses(c echo.Context) error {
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	status := read_topology_deployment_status.NewReadDeploymentStatusesGroup()
	err := status.ReadStatus(ctx, t.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("TopologyDeploymentStatusRequest: ReadDeployedTopologyStatuses")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, status.Slice)
}

type ClusterDeploymentStatusRequest struct {
	CloudCtxNsID string `json:"cloudCtxNsID"`
}

func (t *ClusterDeploymentStatusRequest) ReadDeployedClusterTopologies(c echo.Context) error {
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	id, err := strconv.Atoi(t.CloudCtxNsID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ClusterDeploymentInfoRequest: ReadDeployedClusterTopologies")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	status := read_topology_deployment_status.NewReadDeploymentStatusesGroup()
	err = status.ReadLatestDeployedClusterTopologies(ctx, id, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ClusterDeploymentInfoRequest: ReadDeployedClusterTopologies")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, status.Slice)
}
