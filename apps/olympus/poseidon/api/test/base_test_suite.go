package poseidon_base_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

type PoseidonBaseTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D        test_suites.DatastoresTestSuite
	Ts       chronos.Chronos
	Endpoint string

	zeus_client.ZeusClient
}

func (t *PoseidonBaseTestSuite) SetupTest() {
	t.InitLocalConfigs()
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
	t.ZeusClient = zeus_client.NewZeusClient("http://localhost:9010", t.Tc.ProductionLocalTemporalBearerToken)
}

func (t *PoseidonBaseTestSuite) PostRequest(ctx context.Context, endpoint string, payload any) ([]byte, error) {
	t.PrintReqJson(payload)

	resp, err := t.R().
		SetBody(payload).
		Post(endpoint)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("HydraBaseTestSuite: PostRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return nil, err
	}
	t.PrintRespJson(resp.Body())
	return resp.Body(), err
}
