package artemis_trading_auxiliary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type ArtemisAuxillaryTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
	acc            accounts.Account
	acc2           accounts.Account
	goerliWeb3User web3_client.Web3Client
	goerliNode     string
	nonceOffset    int
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
