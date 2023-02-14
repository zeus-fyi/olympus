package hydra_eth2_web3signer

import (
	"github.com/status-im/keycard-go/hexutils"
	bls_signer "github.com/zeus-fyi/zeus/pkg/crypto/bls"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"testing"

	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

type HydraSigningRequestsTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *HydraSigningRequestsTestSuite) TestMockWeb3SignerConsensys() {
	pubkey := "0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"
	att := consensys_eth2_openapi.AttestationSigning{
		Type: ATTESTATION,
		ForkInfo: consensys_eth2_openapi.SigningForkInfo{
			Fork: consensys_eth2_openapi.Fork{
				PreviousVersion: "0x3000101b",
				CurrentVersion:  "0x4000101b",
				Epoch:           "5",
			},
			GenesisValidatorsRoot: "0xd9a65b284eebd5fe2fc797b9793e1a019073b8de3ffeb83056167b4878fb1557",
		},
		SigningRoot: "0x133f9cee5a36d56ca3085db561375b7f668b12f8d8e971aac8578557ca37635f",
		Attestation: consensys_eth2_openapi.AttestationData{
			Slot:            "29270",
			Index:           "0",
			BeaconBlockRoot: "0x830b0c910c98a5669092f67dc30bc29926023a67f5d0792d427850867cd4418c",
			Source: consensys_eth2_openapi.Checkpoint{
				Epoch: "913",
				Root:  "0xc94343e972eb9373ae894dc573aceb6f0660f791edb2d71797f93ef0a47cda38",
			},
			Target: consensys_eth2_openapi.Checkpoint{
				Epoch: "914",
				Root:  "0x0920af50a400c99dd73941a5a0caf2ca36daa9aea34889aaaf1c5abc4c0a67f6",
			},
		},
	}
	resp, err := t.PostRequest(ctx, Eth2SignRequestWithPubkey(pubkey), att)
	t.Require().Nil(err)
	expSig := string(resp)

	keyTwo := bls_signer.NewEthSignerBLSFromExistingKey(t.Tc.TestEthKeyTwoBLS)
	expPubkey := "8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"
	t.Require().Equal(expPubkey, keyTwo.PublicKeyString())

	trimmedHex := strings_filter.Trim0xPrefix(att.SigningRoot)
	signingRoot := hexutils.HexToBytes(trimmedHex)

	sig := keyTwo.Sign(signingRoot)
	signedLocal := "0x" + bls_signer.ConvertBytesToString(sig.Marshal())

	t.Require().Equal(expSig, signedLocal)
}

func TestHydraSigningRequestsTestSuite(t *testing.T) {
	suite.Run(t, new(HydraSigningRequestsTestSuite))
}
