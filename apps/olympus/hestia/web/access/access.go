package hestia_access_keygen

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	return c.JSON(http.StatusOK, resp)
}
