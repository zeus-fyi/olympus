package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	flows_admin "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/flows/admin"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type GetAdminFlowStatsRequest struct{}

func GetAdminFlowStatsHandler(c echo.Context) error {
	request := new(GetAdminFlowStatsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetAdminFlowStatsRequest(c) // Pass the ID to the method
}

func (e *GetAdminFlowStatsRequest) GetAdminFlowStatsRequest(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	evs, err := flows_admin.SelectUserFlowStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, evs)
}
