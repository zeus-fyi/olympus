package ethereum_slashing_protection_watermarking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	"k8s.io/apimachinery/pkg/util/rand"
)

type WatermarkAttestationsWeb3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

var ctx = context.Background()

func (t *WatermarkAttestationsWeb3SignerTestSuite) TestAttestationWatermarkerDoubleVote() {
	Network = rand.String(10)
	pubkey := "0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"
	data1Att := consensys_eth2_openapi.AttestationData{
		Slot:            "1",
		Index:           "3",
		BeaconBlockRoot: "0x8258f4ec23d5e11",
		Source: consensys_eth2_openapi.Checkpoint{
			Epoch: "2",
			Root:  "0x8258f4ec23d5e113f2b62caa40d77d5",
		},
		Target: consensys_eth2_openapi.Checkpoint{
			Epoch: "3",
			Root:  "0x8258f4ec23d5e113f2b62caa40d77d55ga",
		},
	}
	err := WatermarkAttestation(ctx, pubkey, data1Att)
	t.Require().NoError(err)

	// (data_1 == data_2 and data_1.target.epoch == data_2.target.epoch)
	err = WatermarkAttestation(ctx, pubkey, data1Att)
	t.Require().NoError(err)

	// invalid attestation new target < prev target

	dataInvalidTarget := data1Att
	dataInvalidTarget.Source.Epoch = "1"
	dataInvalidTarget.Target.Epoch = "2"
	err = WatermarkAttestation(ctx, pubkey, dataInvalidTarget)
	t.Require().Error(err)

	data2Att := consensys_eth2_openapi.AttestationData{
		Slot:            "1",
		Index:           "6",
		BeaconBlockRoot: "0x8258f2caa40d77d52c2ad9dfd",
		Source: consensys_eth2_openapi.Checkpoint{
			Epoch: "2",
			Root:  "0x8258f4ec23d5e113f2b62caa40d77d5",
		},
		Target: consensys_eth2_openapi.Checkpoint{
			Epoch: "3",
			Root:  "0x8258f4ec23d5e113f2b62caa40d77d55ga",
		},
	}
	//	(data_1 != data_2 and data_1.target.epoch == data_2.target.epoch)
	t.Require().True(data1Att.Target.Epoch == data2Att.Target.Epoch)
	err = WatermarkAttestation(ctx, pubkey, data2Att)
	t.Require().Error(err)

	// valid attestation
	data1Att.Source.Epoch = "4"
	data1Att.Target.Epoch = "5"
	err = WatermarkAttestation(ctx, pubkey, data1Att)
	t.Require().NoError(err)

	// invalid attestation source > target
	data1Att.Source.Epoch = "10"
	data1Att.Target.Epoch = "9"
	err = WatermarkAttestation(ctx, pubkey, data1Att)
	t.Require().Error(err)

	// invalid attestation (data_1 != data_2 and data_1.target.epoch == data_2.target.epoch)
	err = WatermarkAttestation(ctx, pubkey, data2Att)
	t.Require().Error(err)

}

func (t *WatermarkAttestationsWeb3SignerTestSuite) TestAttestationWatermarkerSurroundVote() {
	Network = rand.String(10)
	pubkey := "0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"
	data1Att := consensys_eth2_openapi.AttestationData{
		Slot:            "1",
		Index:           "",
		BeaconBlockRoot: "",
		Source: consensys_eth2_openapi.Checkpoint{
			Epoch: "2",
			Root:  "",
		},
		Target: consensys_eth2_openapi.Checkpoint{
			Epoch: "3",
			Root:  "",
		},
	}
	sourceEpochData1, targetEpochData1, err := ConvertAttSourceTargetsToInt(ctx, data1Att)
	t.Require().NoError(err)
	err = WatermarkAttestation(ctx, pubkey, data1Att)
	t.Require().NoError(err)

	data2Att := consensys_eth2_openapi.AttestationData{
		Slot:            "1",
		Index:           "",
		BeaconBlockRoot: "",
		Source: consensys_eth2_openapi.Checkpoint{
			Epoch: "3",
			Root:  "",
		},
		Target: consensys_eth2_openapi.Checkpoint{
			Epoch: "2",
			Root:  "",
		},
	}
	sourceEpochData2, targetEpochData2, err := ConvertAttSourceTargetsToInt(ctx, data2Att)
	t.Require().NoError(err)

	isSurrounded := IsSurroundVote(ctx, pubkey, sourceEpochData1, targetEpochData1, sourceEpochData2, targetEpochData2)
	t.Require().True(isSurrounded)
	err = WatermarkAttestation(ctx, pubkey, data2Att)
	t.Require().Error(err)
}

func (t *WatermarkAttestationsWeb3SignerTestSuite) TestFarFutureSigning() {
	// TODO check source and target epochs
}

func TestWatermarkAttestationsWeb3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(WatermarkAttestationsWeb3SignerTestSuite))
}
