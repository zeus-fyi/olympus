package hydra_eth2_web3signer

import (
	"fmt"
	"github.com/status-im/keycard-go/hexutils"
	bls_signer "github.com/zeus-fyi/zeus/pkg/crypto/bls"
	"testing"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
)

type HydraEth2Web3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

// 0xa572a831f2b20ae2016469efe38c0905423ba18105ff46c8dfd2bfa674f36f9e5fcd2588b1d61bf683ff36e983017e971658a961b37a45ea3e70fed99c953b6de4125938901f548e44fd3f5bc44543c88c0b02cf8c08dfdd98dac38ba0d3f395
func (t *HydraEth2Web3SignerTestSuite) TestSigVerify() {
	b := hexutils.HexToBytes("8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08")
	strHex := hexutils.BytesToHex(b)

	t.Require().NotEmpty(strHex)

	err := bls_signer.InitEthBLS()
	t.Require().Nil(err)

	keyOne := bls_signer.NewEthSignerBLSFromExistingKey(t.Tc.TestEthKeyOneBLS)
	expPubkey := "8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288"
	t.Require().Equal(expPubkey, keyOne.PublicKeyString())

	// key 0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08
	//sig := "0x90b9b9d28e47832f888667c7132e3b55baa6f9ff985a4aded492bdf2a27352f67377b6ebb0f3cb6a1e19933f7be9fce019c63226c5282a947581f53fa0832db4602589379c81dd84fb6fcd6e2738102fb5d384b733cad3676549983e764dd6c4"

	// sr 0x4bb4ba9718401008c0a92d3b17b050209151fd7300357d328e5ae6a973f22356

	keyTwo := bls_signer.NewEthSignerBLSFromExistingKey(t.Tc.TestEthKeyTwoBLS)
	expPubkey = "8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"
	t.Require().Equal(expPubkey, keyTwo.PublicKeyString())

	signingRoot := "0x4bb4ba9718401008c0a92d3b17b050209151fd7300357d328e5ae6a973f22356"
	signature := keyTwo.Sign([]byte(signingRoot))

	verified := signature.Verify([]byte(signingRoot), keyTwo.PublicKey())
	t.Require().True(verified)

	expSig := "0x90b9b9d28e47832f888667c7132e3b55baa6f9ff985a4aded492bdf2a27352f67377b6ebb0f3cb6a1e19933f7be9fce019c63226c5282a947581f53fa0832db4602589379c81dd84fb6fcd6e2738102fb5d384b733cad3676549983e764dd6c4"

	returnedSigStr := bls_signer.ConvertBytesToString(signature.Marshal())
	t.Assert().Equal(expSig, returnedSigStr)
}

func (t *HydraEth2Web3SignerTestSuite) TestQueueResponse() {
	newUUID := uuid.New()
	sr := SignRequest{
		UUID:        newUUID,
		Pubkey:      "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		Type:        AGGREGATION_SLOT,
		SigningRoot: "gadgdsgdsg",
	}
	mockSig := "dsadfgdasf"
	SignatureResponsesCache.Set(newUUID.String(), mockSig, cache.DefaultExpiration)
	resp := ReturnSignedMessage(ctx, sr)
	t.Assert().Equal(mockSig, resp.Signature)
}

func (t *HydraEth2Web3SignerTestSuite) TestAggregationSlot() {
	ags := t.GenerateMockAggregationSlotSigningRequest()
	pubkey := "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288"
	ags.Type = AGGREGATION_SLOT
	resp, err := t.PostRequest(ctx, Eth2SignRequestWithPubkey(pubkey), ags)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func (t *HydraEth2Web3SignerTestSuite) TestAttestation() {
	t.InitLocalConfigs()
	t.E.POST(Eth2SignRoute, HydraEth2SignRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9000")
	}()

	<-start
	defer t.E.Shutdown(ctx)

	att := t.GenerateMockAttestationSigningRequest()
	pubkey := "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288"
	resp, err := t.PostRequest(ctx, Eth2SignRequestWithPubkey(pubkey), att)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func Eth2SignRequestWithPubkey(pubkey string) string {
	return fmt.Sprintf("/api/v1/eth2/sign/%s", pubkey)
}

func TestHydraEth2Web3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(HydraEth2Web3SignerTestSuite))
}
