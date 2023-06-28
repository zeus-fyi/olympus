package tyche_base_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type TycheBaseTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D        test_suites.DatastoresTestSuite
	Ts       chronos.Chronos
	Endpoint string

	zeus_client.ZeusClient
}

func (t *TycheBaseTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.E = echo.New()
	//t.Eg = t.E.Group("/")
	//t.Eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	//	AuthScheme: "Bearer",
	//	Validator: func(token string, c echo.Context) (bool, error) {
	//		ctx := context.Background()
	//		key, err := auth.VerifyInternalBearerToken(ctx, token)
	//		if err != nil {
	//			log.Err(err).Msg("InitV1InternalRoutes")
	//			return false, c.JSON(http.StatusInternalServerError, nil)
	//		}
	//		ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
	//		c.Set("orgUser", ou)
	//		c.Set("bearer", key.PublicKey)
	//		return key.PublicKeyVerified, err
	//	},
	//}))
	t.ZeusClient = zeus_client.NewZeusClient("http://localhost:9000", t.Tc.ProductionLocalTemporalBearerToken)
}

func (t *TycheBaseTestSuite) PostRequest(ctx context.Context, endpoint string, payload any) ([]byte, error) {
	t.PrintReqJson(payload)

	resp, err := t.R().
		SetBody(payload).
		Post(endpoint)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("TycheBaseTestSuite: PostRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return nil, err
	}
	t.PrintRespJson(resp.Body())
	return resp.Body(), err
}
