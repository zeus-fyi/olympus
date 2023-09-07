package hestia_iris_dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
	"golang.org/x/net/context"
)

func TutorialToggleRequestHandler(c echo.Context) error {
	request := new(TutorialToggleRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ToggleTutorialSetting(c)
}

type TutorialToggleRequest struct{}

func (t *TutorialToggleRequest) ToggleTutorialSetting(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Msg("failed to get orgUser")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tutorialOn, err := hestia_quicknode_models.ToggleTutorialSetting(context.Background(), ou.OrgID)
	if err != nil {
		log.Err(err).Msg("failed to toggle tutorial setting")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, tutorialOn)
}
