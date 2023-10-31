package v3_iris_internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type IrisV3TestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *IrisV3TestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	s.E = echo.New()

	iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
}

func (s *IrisV3TestSuite) TestThing() {
	InitV3InternalRoutes1(s.E)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = s.E.Start(":9011")
	}()

	token, err := iris_redis.IrisRedisClient.SetInternalAuthCache(ctx, org_users.OrgUser{
		OrgUsers: hestia_autogen_bases.OrgUsers{
			UserID: s.Tc.ProductionLocalTemporalUserID,
			OrgID:  s.Tc.ProductionLocalTemporalOrgID,
		},
	}, "sessionID", "enterprise", "sessionID")
	s.Require().NoError(err)

	ous, plan, err := iris_redis.IrisRedisClient.GetInternalAuthCacheIfExists(ctx, token)
	s.Require().NoError(err)
	s.Require().NotNil(ous)
	s.Require().NotNil(plan)
	s.Assert().Equal(s.Tc.ProductionLocalTemporalOrgID, ous.OrgID)
	s.Assert().Equal(s.Tc.ProductionLocalTemporalUserID, ous.UserID)
	s.Assert().Equal("enterprise", plan)

	groupName := "ethereum-mainnet"
	irisClient := resty_base.GetBaseRestyClient("http://localhost:9011", token)
	irisClient.Header.Set("Content-Type", "application/json")
	irisClient.Header.Set(iris_programmable_proxy.RouteGroupHeader, groupName)
	irisClient.Header.Set(iris_programmable_proxy_v1_beta.LoadBalancingStrategy, iris_programmable_proxy_v1_beta.Adaptive)
	payload := `{
		"jsonrpc": "2.0",
		"procedure": "eth_maxBlockAggReduce",
		"method": "eth_getBlockByNumber",
		"params": ["latest", true],
		"id": 1
	}`
	resp, err := irisClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/v3/internal/router")

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotEqual(http.StatusNotFound, resp.StatusCode())
}

func TestIrisV3TestSuite(t *testing.T) {
	suite.Run(t, new(IrisV3TestSuite))
}
