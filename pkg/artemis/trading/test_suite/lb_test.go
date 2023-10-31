package artemis_trading_test_suite

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
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

func TestLbEnvTradingTestSuite(t *testing.T) {
	suite.Run(t, new(LbEnvTradingTestSuite))
}
