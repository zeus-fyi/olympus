package create_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/cluster"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
)

type TopologyCreateClusterRequest struct {
	ClusterName string   `json:"name"`
	Bases       []string `json:"bases,omitempty"`
}

type TopologyCreateClassResponse struct {
	ClusterName string `json:"name,omitempty"`
	ClassID     int    `json:"classID"`
	Status      string `json:"status,omitempty"`
}

func (t *TopologyCreateClusterRequest) CreateTopologyClusterClass(c echo.Context) error {
	// from auth lookup
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	cc := create_clusters.NewClusterClassTopologyTypeWithBases(ou.OrgID, t.ClusterName, t.Bases)
	err := create_systems.InsertSystem(ctx, &cc.Systems)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("CreateTopologyClusterClass: InsertSystem")
		return c.JSON(http.StatusInternalServerError, err)
	}
	resp := TopologyCreateClassResponse{
		ClusterName: t.ClusterName,
		ClassID:     cc.TopologySystemComponentID,
	}
	return c.JSON(http.StatusOK, resp)
}
