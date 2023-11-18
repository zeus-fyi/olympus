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
	ClusterClassName  string   `json:"clusterClassName"`
	ComponentBaseName string   `json:"componentBaseName,omitempty"`
	SkeletonBaseNames []string `json:"skeletonBaseNames,omitempty"`
}

type TopologyCreateOrAddComponentBasesToClassesRequest struct {
	ClusterClassName   string   `json:"clusterClassName,omitempty"`
	ComponentBaseNames []string `json:"componentBaseNames,omitempty"`
}

type TopologyCreateClassResponse struct {
	ClassID          int    `json:"classID"`
	ClusterClassName string `json:"clusterClassName,omitempty"`
	Status           string `json:"status,omitempty"`
}

func (t *TopologyCreateOrAddBasesToClassesRequest) CreateTopologyClusterClass(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()

	if len(t.ClusterClassName) <= 0 {
		return c.JSON(http.StatusBadRequest, "ClusterClassName is required")
	}
	cc := create_clusters.NewClusterClassTopologyTypeWithBases(ou.OrgID, t.ClusterClassName, t.SkeletonBaseNames)
	err := create_systems.InsertSystem(ctx, &cc.Systems)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("CreateTopologyClusterClass: InsertSystem")
		return c.JSON(http.StatusInternalServerError, err)
	}
	resp := TopologyCreateClassResponse{
		ClusterClassName: t.ClusterClassName,
		ClassID:          cc.TopologySystemComponentID,
	}
	return c.JSON(http.StatusOK, resp)
}

func (t *TopologyCreateOrAddComponentBasesToClassesRequest) AddComponentBasesToTopologyClusterClass(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()
	if len(t.ClusterClassName) <= 0 {
		return c.JSON(http.StatusBadRequest, "ClusterClassName is required")
	}
	bs := make([]bases.Base, len(t.ComponentBaseNames))
	for i, b := range t.ComponentBaseNames {
		bs[i] = bases.NewBaseClassTopologyInsert(ou.OrgID, b)
	}
	err := create_bases.InsertBases(ctx, ou.OrgID, t.ClusterClassName, bs)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TopologyAddBasesToClusterRequest: AddBasesToTopologyClusterClass")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := TopologyCreateClassResponse{
		ClusterClassName: t.ClusterClassName,
	}
	return c.JSON(http.StatusOK, resp)
}
