package artemis_trading_test_suite

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	artemis_mev_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type LbEnvTradingTestSuite struct {
	ArtemisTradingTestSuite
}

func (s *LbEnvTradingTestSuite) TestHypnosLocal() {
	pkHexString2 := s.Tc.ArtemisGoerliEcdsaKey
	secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
	s.Assert().Nil(err)

	iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
	token, err := iris_redis.IrisRedisClient.SetInternalAuthCache(ctx, org_users.OrgUser{
		OrgUsers: hestia_autogen_bases.OrgUsers{
			UserID: s.Tc.ProductionLocalTemporalUserID,
			OrgID:  s.Tc.ProductionLocalTemporalOrgID,
		},
	}, s.Tc.ProductionLocalTemporalBearerToken, "enterprise", "test")
	s.Require().NoError(err)
	fmt.Println(token)
	local := "http://localhost:8080/v1/router"
	wa := web3_client.NewWeb3Client(local, secondAccount)
	wa.AddDefaultEthereumMainnetTableHeader()
	wa.AddSessionLockHeader("test")
	wa.IsAnvilNode = true
	wa.Dial()
	defer wa.Close()
	//origInfo, err := wa.GetNodeMetadata(ctx)
	//s.NoError(err)
	//s.NotEmpty(origInfo)
	//fmt.Println(origInfo.ForkConfig.ForkUrl)

	wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	rpcNew := "http://localhost:8888/node"
	err = wa.ResetNetwork(ctx, rpcNew, 0)
	s.NoError(err)

	origInfo, err := wa.GetNodeMetadata(ctx)
	s.NoError(err)
	s.NotEmpty(origInfo)
	s.Assert().Equal(rpcNew, origInfo.ForkConfig.ForkUrl)
	fmt.Println(origInfo.ForkConfig.ForkUrl)
}

func (s *LbEnvTradingTestSuite) TestEndSessionID() {
	rpcNew := "https://iris.zeus.fyi/v1/router"

	wa := web3_client.NewWeb3ClientFakeSigner(rpcNew)
	wa.AddDefaultEthereumMainnetTableHeader()
	sessionID := "672d8815-e6f2-4040-bbd9-d60337418d64"
	wa.AddSessionLockHeader(sessionID)
	wa.AddEndSessionLockHeader(sessionID)
	wa.IsAnvilNode = true
	wa.Dial()
	defer wa.Close()

	wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)

	origInfo, err := wa.GetNodeMetadata(ctx)
	s.NoError(err)
	s.NotEmpty(origInfo)
	fmt.Println(origInfo.ForkConfig.ForkUrl)
}

func (s *LbEnvTradingTestSuite) TestFork() {
	rpcNew := "https://iris.zeus.fyi/v1/router"
	wa := web3_client.NewWeb3ClientFakeSigner(rpcNew)
	wa.AddDefaultEthereumMainnetTableHeader()
	sessionID := uuid.New().String()
	artemis_orchestration_auth.Bearer = s.Tc.ProductionLocalTemporalBearerToken

	defer func(sessionID string) {
		artemis_mev_transcations.EndServerlessEnvironment(sessionID)
	}(sessionID)
	wa.AddSessionLockHeader(sessionID)
	wa.IsAnvilNode = true
	wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wa.Dial()
	defer wa.Close()

	origInfo, err := wa.GetNodeMetadata(ctx)
	s.NoError(err)
	s.NotEmpty(origInfo)
	fmt.Println(origInfo.ForkConfig.ForkUrl)

	/*
		rpcNew := "https://iris.zeus.fyi/v1/router"

		wa := web3_client.NewWeb3ClientFakeSigner(rpcNew)
		wa.AddDefaultEthereumMainnetTableHeader()
		sessionID := "672d8815-e6f2-4040-bbd9-d60337418d64"
		wa.AddSessionLockHeader(sessionID)
		wa.AddEndSessionLockHeader(sessionID)
		wa.IsAnvilNode = true
		wa.Dial()
		defer wa.Close()

		wa.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)

		origInfo, err := wa.GetNodeMetadata(ctx)
		s.NoError(err)
		s.NotEmpty(origInfo)
		fmt.Println(origInfo.ForkConfig.ForkUrl)

		pkHexString2 := s.Tc.LocalEcsdaTestPkey2
		secondAccount, err := accounts.ParsePrivateKey(pkHexString2)
		user := web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, secondAccount)
		uswap := web3_client.InitUniswapClient(ctx, user)

		uswap.Web3Client.Dial()
		uswap.Web3Client.AddDefaultEthereumMainnetTableHeader()
		defer uswap.Web3Client.Close()
		err = uswap.Web3Client.ResetNetwork(ctx, "https://localhost:8545", 100000)
		s.NoError(err)

	*/
}

func TestLbEnvTradingTestSuite(t *testing.T) {
	suite.Run(t, new(LbEnvTradingTestSuite))
}
