package read_infra

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

type TopologyReadPrivateAppsRequest struct {
}

func (t *TopologyReadPrivateAppsRequest) ListPrivateAppsRequest(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	apps, err := read_topology.SelectOrgApps(ctx, ou.OrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, apps)
}

func (t *TopologyReadPrivateAppsRequest) GetPrivateAppDetailsRequest(c echo.Context) error {
	appID := c.Param("id")
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	id, err := strconv.Atoi(appID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("GetPrivateAppDetailsRequest: SelectOrgApps")
		return c.JSON(http.StatusBadRequest, nil)
	}
	apps, err := read_topology.SelectAppTopologyByID(ctx, ou.OrgID, id)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := zeus_templates.ClusterPreviewWorkloadsOlympus{
		ClusterName:    apps.ClusterClassName,
		ComponentBases: make(map[string]map[string]any),
	}

	for cbName, cb := range apps.ComponentBases {
		uiSbs := make(map[string]any)
		resp.ComponentBases[cbName] = make(map[string]any)
		for sbName, sb := range cb.SkeletonBases {
			tr := read_topology.NewInfraTopologyReader()
			tr.TopologyID = sb.TopologyID
			// from auth lookup
			tr.OrgID = ou.OrgID
			tr.UserID = ou.UserID
			err = tr.SelectTopologyForOrg(ctx)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectTopology")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			nk := tr.GetTopologyBaseInfraWorkload()
			uiSbs[sbName] = nk
		}
		resp.ComponentBases[cbName] = uiSbs
	}
	return c.JSON(http.StatusOK, resp)
}
