package artemis_realtime_trading

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type ArtemisRealTimeTradingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	MainnetWeb3User    web3_client.Web3Client
	ProxiedMainnetUser web3_client.Web3Client
	at                 ActiveTrading
}

func (s *ArtemisRealTimeTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	artemis_trading_cache.InitTokenFilter(ctx)

	artemis_test_cache.InitLiveTestNetwork(s.Tc.QuikNodeURLS.TestRoute)
	pkHexString := s.Tc.LocalEcsdaTestPkey

	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)

	pkHexString2 := s.Tc.LocalEcsdaTestPkey2
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)
	s.MainnetWeb3User = web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
	s.ProxiedMainnetUser = web3_client.NewWeb3Client(artemis_trading_constants.IrisAnvilRoute, secondAccount)
	s.ProxiedMainnetUser.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)

	uni := web3_client.InitUniswapClient(ctx, s.ProxiedMainnetUser)
	uni.PrintOn = true
	uni.PrintLocal = true
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	uni.DebugPrint = true

	s.at = NewActiveTradingDebugger(&uni)
	artemis_trading_cache.InitTokenFilter(ctx)
}

func TestArtemisRealTimeTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisRealTimeTradingTestSuite))
}
