package pods

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ExternalApiPodsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.Use(PodsCloudCtxNsMiddleware)
	e.POST("/pods", HandlePodActionRequest)
	return e
}

func PodsCloudCtxNsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Info().Msg("PodsCloudCtxNsMiddleware")
		ctx := context.Background()
		ou := c.Get("orgUser").(org_users.OrgUser)
		request := new(PodActionRequest)
		if err := c.Bind(request); err != nil {
			return err
		}

		cctxID := c.Request().Header.Get("CloudCtxNsID")
		if len(cctxID) == 0 {
			authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
			if authed != true {
				return c.JSON(http.StatusUnauthorized, err)
			}
			c.Set("PodActionRequest", request)
			return next(c)
		}
		cID, err := strconv.Atoi(cctxID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, ou.OrgID, cID)
		if authed != true {
			return c.JSON(http.StatusUnauthorized, nil)
		}
		request.CloudCtxNs = cctx
		c.Set("PodActionRequest", request)
		return next(c)
	}
}
