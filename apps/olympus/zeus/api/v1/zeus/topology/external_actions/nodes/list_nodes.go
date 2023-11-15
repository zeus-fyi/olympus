package nodes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ListNodesRequest(c echo.Context, request *ActionRequest) error {
	ctx := context.Background()
	log.Info().Msg("ListNodesRequest")
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Msg("ListNodesRequest, no orgUser found")
		return c.JSON(http.StatusBadRequest, nil)
	}
	label := fmt.Sprintf("org=%d", ou.OrgID)
	for k, v := range request.Labels {
		label += fmt.Sprintf(",%s=%s", k, v)
	}
	nl, err := zeus.K8Util.GetNodesAuditByLabel(ctx, request.CloudCtxNs, label)
	if err != nil {
		log.Err(err).Msg("ListNodesRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nl)
}
