package artemis_rawdawg_contract

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
	"github.com/zeus-fyi/zeus/examples/adaptive_rpc_load_balancer/smart_contract_library"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

var ctx = context.Background()

type ArtemisTradingContractsTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
}

func (s *ArtemisTradingContractsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

var LoadBalancerAddress = "https://iris.zeus.fyi/v1/router"

func CreateLocalUser(ctx context.Context, bearer, sessionID string) web3_actions.Web3Actions {
	acc, err := accounts.CreateAccount()
	if err != nil {
		panic(err)
	}
	w3a := web3_actions.NewWeb3ActionsClientWithAccount(LoadBalancerAddress, acc)
	w3a.AddAnvilSessionLockHeader(sessionID)
	w3a.AddBearerToken(bearer)
	nvB := (*hexutil.Big)(smart_contract_library.EtherMultiple(10000))
	w3a.Dial()
	defer w3a.Close()
	err = w3a.SetBalance(ctx, w3a.Address().String(), *nvB)
	if err != nil {
		panic(err)
	}
	return w3a
}

func TestArtemisTradingContractsTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradingContractsTestSuite))
}
