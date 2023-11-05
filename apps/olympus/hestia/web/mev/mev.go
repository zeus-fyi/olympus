package hestia_mev

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
)

type MevRequest struct {
}

func MevRequestHandler(c echo.Context) error {
	request := new(MevRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetDashboardInfo(c)
}

var (
	PromqlProxy = "http://promql.promql-edc89f30.svc.cluster.local"
)

func (r *MevRequest) GetDashboardInfo(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if ou.OrgID != hestia_login.TemporalOrgID {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ctx := context.Background()
	bundles, err := artemis_reporting.GetBundleSubmissionHistory(ctx, 0, 1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	di := bundles.GetDashboardInfo()
	rc := resty.New()
	_, err = rc.R().SetResult(&di.TopKTokens).Get(PromqlProxy + "/v1/promql/top/tokens")
	if err != nil {
		log.Err(err).Msg("failed to get top tokens")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, di)
}
