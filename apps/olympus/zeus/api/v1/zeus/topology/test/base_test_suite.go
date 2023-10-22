package test

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
)

var Kns zeus_common_types.CloudCtxNs

var TestOrgUser = org_users.OrgUser{autogen_bases.OrgUsers{
	OrgID:  1667452524363177528,
	UserID: 1667452524356256466,
}}

var TestTopologyID = 7140168037686545724

type TopologyActionRequestTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D        test_suites.DatastoresTestSuite
	Ts       chronos.Chronos
	Endpoint string

	zeus_client.ZeusClient
}

type TestResponse struct {
	Logs []byte
	pods v1.PodList
}

func (t *TopologyActionRequestTestSuite) SetupTest() {
	t.K.ConnectToK8s()
	t.InitLocalConfigs()

	zeus.K8Util = t.K
	t.D.PGTest.SetupPGConn()
	t.D.PG = t.D.PGTest.Pg
	t.E = echo.New()

	eg := t.E.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerToken(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))

	t.Eg = eg
	t.ZeusClient = zeus_client.NewZeusClient("http://localhost:9010", t.Tc.LocalBearerToken)
}

func (t *TopologyActionRequestTestSuite) AddEndpointHandler(h echo.HandlerFunc) {
	t.E.POST(t.Endpoint, h)
}
