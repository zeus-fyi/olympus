package nodes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func ListNodesRequest(c echo.Context, request *ActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Info().Msg("ListNodesRequest")
	ou := c.Get("orgUser").(org_users.OrgUser)
	label := fmt.Sprintf("org=%d", ou.OrgID)
	for k, v := range request.Labels {
		label += fmt.Sprintf(",%s=%s", k, v)
	}
	// TODO update later
	cloudCtx := zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "do-nyc1-do-nyc1-zeus-demo",
		Namespace:     "",
		Env:           "",
	}
	nl, err := zeus.K8Util.GetNodesAuditByLabel(ctx, cloudCtx, label)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ListNodesRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nl)
}
