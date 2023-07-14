package iris_base_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type IrisBaseTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D               test_suites.DatastoresTestSuite
	Ts              chronos.Chronos
	Endpoint        string
	MainnetWeb3User web3_client.Web3Client

	zeus_client.ZeusClient
}

func (t *IrisBaseTestSuite) SetupTest() {
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
	t.ZeusClient = zeus_client.NewZeusClient("http://localhost:8080", t.Tc.ProductionLocalTemporalBearerToken)
}

func (t *IrisBaseTestSuite) PostRequest(ctx context.Context, endpoint string, payload any) ([]byte, error) {
	t.PrintReqJson(payload)

	resp, err := t.ZeusClient.R().
		SetBody(payload).
		Post(endpoint)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("IrisBaseTestSuite: PostRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return nil, err
	}
	t.PrintRespJson(resp.Body())
	return resp.Body(), err
}
