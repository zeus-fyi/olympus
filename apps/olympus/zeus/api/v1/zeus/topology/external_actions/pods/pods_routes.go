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
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
)

func ExternalApiPodsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.Use(PodsCloudCtxNsMiddlewareWrapper(k8Cfg))
	e.POST("/pods", HandlePodActionRequest)
	return e
}

const (
	TemporalOrgID   = 7138983863666903883
	KevinFlowsOrgID = 1710298581127603000
)

func PodsCloudCtxNsMiddlewareWrapper(k8Cfg autok8s_core.K8Util) echo.MiddlewareFunc {
	// Return a function that conforms to Echo's MiddlewareFunc signature
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Return the actual middleware function, utilizing k8Cfg as needed
		return func(c echo.Context) error {
			log.Info().Msg("PodsCloudCtxNsMiddleware")
			ctx := context.Background()
			ou, ok := c.Get("orgUser").(org_users.OrgUser)
			if !ok {
				log.Warn().Msg("PodsCloudCtxNsMiddleware: orgUser not found")
				return c.JSON(http.StatusUnauthorized, nil)
			}
			c.Set("orgUser", ou)
			request := new(zeus_pods_reqs.PodActionRequest)
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
				log.Err(err).Msg("PodsCloudCtxNsMiddleware")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			tmp := ou.OrgID
			if cID == 1715903995547997000 && tmp == KevinFlowsOrgID {
				tmp = TemporalOrgID
			}
			authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, tmp, cID)
			if authed != true {
				log.Warn().Interface("ou", ou).Interface("req", request).Msg("PodsCloudCtxNsMiddleware: IsOrgCloudCtxNsAuthorizedFromID")

				return c.JSON(http.StatusUnauthorized, nil)
			}
			request.CloudCtxNs = cctx
			k, err := zeus.VerifyClusterAuthFromCtxOnlyAndGetKubeCfg(c.Request().Context(), ou, cctx)
			if err != nil {
				log.Warn().Interface("ou", ou).Interface("req", request).Msg("PodsCloudCtxNsMiddleware: IsOrgCloudCtxNsAuthorizedFromID")
				return c.JSON(http.StatusUnauthorized, nil)
			}
			if k != nil {
				c.Set("k8Cfg", *k)
			} else {
				c.Set("k8Cfg", k8Cfg)
			}
			c.Set("PodActionRequest", request)
			return next(c)
		}
	}
}
