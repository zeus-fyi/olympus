package artemis_trading_test_suite

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

type ArtemisTradingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	MainnetWeb3User    web3_client.Web3Client
	ProxiedMainnetUser web3_client.Web3Client
	GoerliWeb3User     web3_client.Web3Client

	LocalIrisAnvilWeb3User web3_client.Web3Client
	IrisAnvilWeb3User      web3_client.Web3Client
}

var ctx = context.Background()

func (s *ArtemisTradingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	artemis_test_cache.InitLiveTestNetwork(s.Tc.QuikNodeURLS.TestRoute)

	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	artemis_trading_cache.InitTokenFilter(ctx)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)

	s.MainnetWeb3User = web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
	wc := web3_client.NewWeb3Client(artemis_trading_constants.IrisExtAnvilRoute, newAccount)
	wc.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wc.AddSessionLockHeader(uuid.New().String())
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	s.ProxiedMainnetUser = wc

	pkHexString2 := s.Tc.ArtemisGoerliEcdsaKey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)
	s.GoerliWeb3User = web3_client.NewWeb3Client(s.Tc.GoerliNodeUrl, secondAccount)
	s.MainnetWeb3User = web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, secondAccount)
}

func (s *ArtemisTradingTestSuite) TestLb() {
	// "http://localhost:8080/v1/router"
	pkHexString2 := s.Tc.ArtemisGoerliEcdsaKey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)

	//wa := web3_client.NewWeb3Client(artemis_trading_constants.IrisExtAnvilRoute, secondAccount)
	//wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	//wa.AddSessionLockHeader("test")
	//wa.IsAnvilNode = true
	//s.IrisAnvilWeb3User = wa

	token, err := iris_redis.IrisRedisClient.SetInternalAuthCache(ctx, org_users.OrgUser{
		OrgUsers: hestia_autogen_bases.OrgUsers{
			UserID: s.Tc.ProductionLocalTemporalUserID,
			OrgID:  s.Tc.ProductionLocalTemporalOrgID,
		},
	}, s.Tc.ProductionLocalTemporalBearerToken, "enterprise", "test")
	s.Require().NoError(err)

	local := "http://localhost:8888"
	wa := web3_client.NewWeb3Client(local, secondAccount)
	wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wa.AddSessionLockHeader(token)
	wa.IsAnvilNode = true
	s.LocalIrisAnvilWeb3User = wa

	wa.AddSessionLockHeader("test")
	wa.AddDefaultEthereumMainnetTableHeader()
	wa.Dial()
	origInfo, err := wa.GetNodeMetadata(ctx)
	s.NoError(err)
	s.NotEmpty(origInfo)
	wa.Close()
	wa.Dial()

	err = wa.ResetNetwork(ctx, "", 0)
	s.NoError(err)
	wa.Close()
	//for i := 0; i < 10; i++ {
	//	wa.Dial()
	//	nodeInfo, rerr := wa.GetNodeMetadata(ctx)
	//	if rerr != nil {
	//		return -1, err
	//	}
	//	wa.Close()
	//	if nodeInfo.ForkConfig.ForkUrl != origInfo.ForkConfig.ForkUrl {
	//		return -1, fmt.Errorf("CheckBlockRxAndNetworkReset: live network fork url %s is not equal to initial fork url %s", nodeInfo.ForkConfig.ForkUrl, nodeInfo.ForkConfig.ForkUrl)
	//	}
	//	if nodeInfo.ForkConfig.ForkBlockNumber != currentBlockNum {
	//		fmt.Println("initForkUrl1", origInfo.ForkConfig.ForkUrl, "CurrentBlockNumber", origInfo.CurrentBlockNumber.ToInt().String(), "ForkBlockNumber", origInfo.ForkConfig.ForkBlockNumber)
	//		fmt.Println("initForkUrl2", nodeInfo.ForkConfig.ForkUrl, "CurrentBlockNumber", nodeInfo.CurrentBlockNumber.ToInt().String(), "ForkBlockNumber", nodeInfo.ForkConfig.ForkBlockNumber)
	//	} else {
	//		return currentBlockNum, nil
	//	}
	//	time.Sleep(100 * time.Millisecond)
	//}
}

func TestArtemisTradingTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradingTestSuite))
}
