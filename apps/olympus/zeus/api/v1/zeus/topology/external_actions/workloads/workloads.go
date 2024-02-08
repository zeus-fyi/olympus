package external_api_workloads

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyCloudCtxNsQueryRequest struct {
	zeus_common_types.CloudCtxNs
}

func (t *TopologyCloudCtxNsQueryRequest) ReadDeployedWorkloads(c echo.Context) error {
	log.Debug().Msg("TopologyCloudCtxNsQueryRequest")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	k, err := zeus.VerifyClusterAuthFromCtxOnlyAndGetKubeCfg(c.Request().Context(), ou, t.CloudCtxNs)
	if err != nil {
		log.Warn().Interface("ou", ou).Interface("req", t).Msg("PodsCloudCtxNsMiddleware: IsOrgCloudCtxNsAuthorizedFromID")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	kCfg := zeus.K8Util
	if k != nil {
		kCfg = *k
	}
	workload, err := kCfg.GetWorkloadAtNamespace(ctx, t.CloudCtxNs)
	if err != nil {
		log.Err(err).Msg("TopologyCloudCtxNsQueryRequest: GetWorkloadAtNamespace")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, workload)
}
