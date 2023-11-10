package destroy_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type ResourceDestroyRequest struct {
	OrgResourceID int `json:"orgResourceID"`
}

func (r *ResourceDestroyRequest) DestroyResource(c echo.Context) error {
	log.Debug().Msg("ResourceDestroyRequest: DestroyResource")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return zeus.ExecuteDestroyResourcesWorkflow(c, ctx, ou, []int{r.OrgResourceID})
}
