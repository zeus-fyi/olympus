package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
)

func DeprovisionRequestHandler(c echo.Context) error {
	request := new(DeprovisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deprovision(c)
}

type DeprovisionRequest struct {
	hestia_quicknode.DeprovisionRequest
}

func (r *DeprovisionRequest) Deprovision(c echo.Context) error {
	ouc := c.Get("orgUser")
	ou, ok := ouc.(org_users.OrgUser)
	if !ok {
		key, err := auth.VerifyQuickNodeToken(context.Background(), r.QuickNodeID)
		if err != nil {
			log.Err(err).Msg("InitV1Routes QuickNode user not found: creating new org")
			err = nil
		}
		ou = org_users.NewOrgUserWithID(key.OrgID, 0)
		c.Set("orgUser", ou)
		c.Set("verified", key.IsVerified())
	}
	dp := r.DeprovisionRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnDeprovisionWorkflow(context.Background(), ou, dp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Status: "error",
			})
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}
