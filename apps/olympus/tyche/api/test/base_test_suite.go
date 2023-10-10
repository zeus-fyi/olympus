package tyche_base_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
)

var ctx = context.Background()

type TycheBaseTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D               test_suites.DatastoresTestSuite
	Ts              chronos.Chronos
	Endpoint        string
	MainnetWeb3User web3_client.Web3Client

	zeus_client.ZeusClient
}

func (t *TycheBaseTestSuite) SetupTest() {
	t.InitLocalConfigs()
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	t.Assert().Nil(err)
	t.MainnetWeb3User = web3_client.NewWeb3Client(t.Tc.QuikNodeLiveNode, newAccount)

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL = t.Tc.QuikNodeLiveNode
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.Account = newAccount
	tyche_metrics.InitTycheMetrics(ctx)
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
