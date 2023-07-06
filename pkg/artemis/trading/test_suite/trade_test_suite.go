package artemis_trading_test_suite

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type ArtemisTradingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	MainnetWeb3User    web3_client.Web3Client
	ProxiedMainnetUser web3_client.Web3Client
}

var ctx = context.Background()

func (s *ArtemisTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	artemis_trading_cache.InitTokenFilter(ctx)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)
	s.MainnetWeb3User = web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
	wc := web3_client.NewWeb3Client(artemis_trading_constants.IrisAnvilRoute, newAccount)
	m := map[string]string{
		"Authorization": "Bearer " + s.Tc.ProductionLocalTemporalBearerToken,
	}
	wc.Headers = m
	wc.AddSessionLockHeader(uuid.New().String())
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	s.ProxiedMainnetUser = wc

	pkHexString2 := s.Tc.LocalEcsdaTestPkey2
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)
	s.MainnetWeb3User = web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, secondAccount)
}

func TestArtemisTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradingTestSuite))
}
