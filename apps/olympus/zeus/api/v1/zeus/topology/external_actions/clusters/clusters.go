package v1_zeus_clusters

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/cloud_ctx_logs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ExternalClusterLogsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.Use(CloudCtxNsMiddlewareWrapper(k8Cfg))
	e.GET("/cluster/logs", ClusterLogsRequestHandler)
	return e
}

func CloudCtxNsMiddlewareWrapper(k8Cfg autok8s_core.K8Util) echo.MiddlewareFunc {
	// Return a function that conforms to Echo's MiddlewareFunc signature
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Return the actual middleware function, utilizing k8Cfg as needed
		return func(c echo.Context) error {
			log.Info().Msg("CloudCtxNsMiddlewareWrapper")
			ctx := context.Background()
			ou, ok := c.Get("orgUser").(org_users.OrgUser)
			if !ok {
				log.Warn().Msg("CloudCtxNsMiddlewareWrapper: orgUser not found")
				return c.JSON(http.StatusUnauthorized, nil)
			}
			c.Set("orgUser", ou)
			request := new(ClusterLogsRequest)
			if err := c.Bind(request); err != nil {
				return err
			}
			cctxID := c.Request().Header.Get("CloudCtxNsID")
			if len(cctxID) == 0 {
				authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, request.CloudCtxNs)
				if authed != true {
					return c.JSON(http.StatusUnauthorized, err)
				}
				c.Set("CloudCtxNsMiddlewareWrapper", request)
				return next(c)
			}
			cID, err := strconv.Atoi(cctxID)
			if err != nil {
				log.Err(err).Msg("CloudCtxNsMiddlewareWrapper")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, ou.OrgID, cID)
			if authed != true {
				log.Warn().Interface("ou", ou).Interface("req", request).Msg("CloudCtxNsMiddlewareWrapper: IsOrgCloudCtxNsAuthorizedFromID")

				return c.JSON(http.StatusUnauthorized, nil)
			}
			request.CloudCtxNs = cctx
			k, err := zeus.VerifyClusterAuthFromCtxOnlyAndGetKubeCfg(c.Request().Context(), ou, cctx)
			if err != nil {
				log.Warn().Interface("ou", ou).Interface("req", request).Msg("CloudCtxNsMiddlewareWrapper: IsOrgCloudCtxNsAuthorizedFromID")
				return c.JSON(http.StatusUnauthorized, nil)
			}
			if k != nil {
				c.Set("k8Cfg", *k)
			} else {
				c.Set("k8Cfg", k8Cfg)
			}
			ccl := cloud_ctx_logs.CloudCtxNsLogs{
				CloudCtxNsID: cID,
				Ou:           ou,
				CloudCtxNs:   request.CloudCtxNs,
			}
			c.Set("CloudCtxNsLogs", ccl)
			c.Set("CloudCtxNsMiddlewareWrapper", request)
			return next(c)
		}
	}
}
