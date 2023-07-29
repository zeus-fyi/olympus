package hestia_quiknode_v1_routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

const (
	QuickNodeTestHeader = "X-QN-TESTING"
	QuickNodeIDHeader   = "x-quicknode-id"
	QuickNodeEndpointID = "x-instance-id"
	QuickNodeChain      = "x-qn-chain"
	QuickNodeNetwork    = "x-qn-network"
)

var (
	QuickNodeUsername  = ""
	QuickNodePassword  = ""
	QuickNodeToken     = ""
	QuickNodeOrgID     = 10
	QuickNodeTestOrgID = 9
)

func InitV1RoutesServices(e *echo.Echo) {
	eg := e.Group("/v1/api")
	eg.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Get headers
		qnTestHeader := c.Request().Header.Get(QuickNodeTestHeader)
		qnIDHeader := c.Request().Header.Get(QuickNodeIDHeader)
		qnEndpointID := c.Request().Header.Get(QuickNodeEndpointID)
		qnChain := c.Request().Header.Get(QuickNodeChain)
		qnNetwork := c.Request().Header.Get(QuickNodeNetwork)
		// Set headers to echo context
		c.Set(QuickNodeTestHeader, qnTestHeader)
		c.Set(QuickNodeIDHeader, qnIDHeader)
		c.Set(QuickNodeEndpointID, qnEndpointID)
		c.Set(QuickNodeChain, qnChain)
		c.Set(QuickNodeNetwork, qnNetwork)
		if len(QuickNodePassword) <= 0 {
			return false, nil
		}

		if password == QuickNodePassword {
			if qnTestHeader == QuickNodeTestHeader {
				c.Set("orgUser", org_users.NewOrgUserWithID(QuickNodeTestOrgID, QuickNodeTestOrgID))
				c.Set("bearer", QuickNodeToken)
				return true, nil
			} else {
				c.Set("orgUser", org_users.NewOrgUserWithID(10, 10))
				c.Set("bearer", QuickNodeToken)
				return true, nil
			}
		}
		return false, nil
	}))

	eg.POST("/provision", ProvisionRequestHandler)
	eg.PUT("/update", UpdateProvisionRequestHandler)
	eg.DELETE("/deprovision", DeprovisionRequestHandler)
	eg.DELETE("/deactivate", DeactivateRequestHandler)
}

type QuickNodeResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
