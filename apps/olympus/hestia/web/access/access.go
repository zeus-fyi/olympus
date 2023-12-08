package hestia_access_keygen

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
)

type AccessRequest struct {
}

func AccessRequestHandler(c echo.Context) error {
	request := new(AccessRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.AuthCheck(c)
}

const (
	TemporalOrgID = 7138983863666903883
	SamsOrgID     = 1701381301753642000
)

func (a *AccessRequest) AuthCheck(c echo.Context) error {
	var resp hestia_login.LoginResponse
	token, ok := c.Get("bearer").(string)
	if ok {
		plan, err := hestia_billing.GetPlan(c.Request().Context(), token)
		if err != nil {
			log.Err(err).Msg("AuthCheck: GetPlan error")
			err = nil
		} else {
			resp.PlanDetailsUsage = &plan
		}
	}
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	isInternal := false
	if ou.OrgID == TemporalOrgID || ou.OrgID == SamsOrgID {
		isInternal = true
	}
	resp.IsInternal = isInternal
	resp.IsBillingSetup = hestia_billing.CheckBillingCache(c.Request().Context(), ou.UserID)
	return c.JSON(http.StatusOK, resp)
}
