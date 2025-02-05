package artemis_trading_auxiliary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type ArtemisAuxillaryTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
	simMainnetTrader AuxiliaryTradingUtils
	at1              AuxiliaryTradingUtils
	at2              AuxiliaryTradingUtils
	atMainnet        AuxiliaryTradingUtils
	acc              accounts.Account
	acc2             accounts.Account
	acc3             accounts.Account
	goerliWeb3User   web3_client.Web3Client
	mainnetWeb3User  web3_client.Web3Client
	mainnetNode      string
	goerliNode       string
	nonceOffset      int
}

var ctx = context.Background()

func (t *ArtemisAuxillaryTestSuite) SetupTest() {
	t.S3SecretsManagerTestSuite.SetupTest()
	t.goerliNode = t.Tc.GoerliNodeUrl
	tc := t.S3SecretsManagerTestSuite.Tc
	athena.AthenaS3Manager = t.S3SecretsManagerTestSuite.S3
	apps.Pg.InitPG(ctx, tc.ProdLocalDbPgconn)
	age := encryption.NewAge(tc.LocalAgePkey, tc.LocalAgePubkey)
	t.acc = initTradingAccount(ctx, age)
	secondAccount, err := accounts.ParsePrivateKey(t.Tc.ArtemisGoerliEcdsaKey)
	t.Assert().Nil(err)
	t.acc2 = *secondAccount
	t.goerliWeb3User = web3_client.NewWeb3Client(t.Tc.GoerliNodeUrl, secondAccount)
	t.mainnetNode = t.Tc.MainnetNodeUrl
	t.mainnetWeb3User = web3_client.NewWeb3Client(t.Tc.MainnetNodeUrl, &t.acc2)
	m := map[string]string{
		"Authorization": "Bearer " + t.Tc.ProductionLocalTemporalBearerToken,
	}
	t.mainnetWeb3User.Headers = m
	network := hestia_req_types.Goerli
	w3a := web3_client.NewWeb3Client(t.goerliNode, &t.acc)
	w3a.Network = network
	w3a2 := web3_client.NewWeb3Client(t.goerliNode, &t.acc2)
	w3a2.Network = network
	t.at1 = InitAuxiliaryTradingUtils(ctx, w3a)
	t.at2 = InitAuxiliaryTradingUtils(ctx, w3a2)

	wc := web3_client.NewWeb3ClientFakeSigner(artemis_trading_constants.IrisAnvilRoute)
	wc.Headers = m
	wc.Network = hestia_req_types.Mainnet
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.DebugPrint = true
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	t.simMainnetTrader = InitAuxiliaryTradingUtilsFromUni(ctx, &uni)
}

// InitTradingAccount pubkey 0x000025e60C7ff32a3470be7FE3ed1666b0E326e2
func initTradingAccount(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-4.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}

// InitTradingAccount pubkey 0x000000641e80A183c8B736141cbE313E136bc8c6
func initTradingAccount2(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-6.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}

func TestArtemisAuxiliaryTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisAuxillaryTestSuite))
}
