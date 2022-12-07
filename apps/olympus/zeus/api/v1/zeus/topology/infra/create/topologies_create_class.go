package create_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	create_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases"
	create_clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/cluster"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
)

type TopologyCreateOrAddBasesToClassesRequest struct {
	ClassName      string   `json:"className"`
	ClassBaseNames []string `json:"classBaseNames,omitempty"`
}

type TopologyCreateClassResponse struct {
	ClusterName string `json:"name,omitempty"`
	ClassID     int    `json:"classID"`
	Status      string `json:"status,omitempty"`
}

func (t *TopologyCreateOrAddBasesToClassesRequest) CreateTopologyClusterClass(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	cc := create_clusters.NewClusterClassTopologyTypeWithBases(ou.OrgID, t.ClassName, t.ClassBaseNames)
	err := create_systems.InsertSystem(ctx, &cc.Systems)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("CreateTopologyClusterClass: InsertSystem")
		return c.JSON(http.StatusInternalServerError, err)
	}
	resp := TopologyCreateClassResponse{
		ClusterName: t.ClassName,
		ClassID:     cc.TopologySystemComponentID,
	}
	return c.JSON(http.StatusOK, resp)
}

func (t *TopologyCreateOrAddBasesToClassesRequest) AddBasesToTopologyClusterClass(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()

	bs := make([]bases.Base, len(t.ClassBaseNames))
	for i, b := range t.ClassBaseNames {
		bs[i] = bases.NewBaseClassTopologyInsert(ou.OrgID, b)
	}
	err := create_bases.InsertBases(ctx, ou.OrgID, t.ClassName, bs)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TopologyAddBasesToClusterRequest: AddBasesToTopologyClusterClass")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
