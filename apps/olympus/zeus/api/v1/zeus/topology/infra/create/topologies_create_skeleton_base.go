package create_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_skeletons "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/skeleton"
)

func (t *TopologyCreateOrAddBasesToClassesRequest) AddSkeletonBaseClassToBase(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	err := create_skeletons.InsertSkeletonBases(ctx, ou.OrgID, t.ClusterClassName, t.ComponentBaseName, t.SkeletonBaseNames)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TopologyCreateSkeletonBasesRequest: AddSkeletonBaseClassToBase")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
