package hydra_eth2_web3signer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

type WatermarkBlockDecodingTestSuite struct {
	suite.Suite
}

var ctx = context.Background()

func (t *WatermarkBlockDecodingTestSuite) TestDecodePhase0() {
	bs := consensys_eth2_openapi.BeaconBlockSigning{}
	blockPhase0 := consensys_eth2_openapi.BlockRequestPhase0{}
	err := faker.FakeData(&blockPhase0)
	t.Require().Nil(err)

	blockPhase0.Block.Slot = "1"
	blockPhase0.Version = PHASE0
	bs.BeaconBlock = blockPhase0
	bs.Type = "BLOCK"

	body, slot, err := DecodeBeaconBlockAndSlot(ctx, bs.BeaconBlock)
	t.Require().Nil(err)
	t.Require().Equal(1, slot)

	typeBody := body.(consensys_eth2_openapi.BlockRequestPhase0)
	rawBytes1, err := json.Marshal(body)
	t.Require().Nil(err)

	rawBytes2, err := json.Marshal(typeBody)
	t.Require().Nil(err)
	t.Require().Equal(rawBytes1, rawBytes2)
}

func (t *WatermarkBlockDecodingTestSuite) TestDecodeAltair() {
	bs := consensys_eth2_openapi.BeaconBlockSigning{}
	blockAltair := consensys_eth2_openapi.BlockRequestAltair{}
	err := faker.FakeData(&blockAltair)
	t.Require().Nil(err)
	blockAltair.Block.Slot = "100000"
	blockAltair.Version = ALTAIR

	bs.Type = "BLOCK_V2"
	bs.BeaconBlock = blockAltair

	body, slot, err := DecodeBeaconBlockAndSlot(ctx, bs.BeaconBlock)
	t.Require().Nil(err)
	t.Require().Equal(100000, slot)

	typeBody := body.(consensys_eth2_openapi.BlockRequestAltair)
	rawBytes1, err := json.Marshal(body)
	t.Require().Nil(err)

	rawBytes2, err := json.Marshal(typeBody)
	t.Require().Nil(err)
	t.Require().Equal(rawBytes1, rawBytes2)
}

func (t *WatermarkBlockDecodingTestSuite) TestDecodeBellatrix() {
	bs := consensys_eth2_openapi.BeaconBlockSigning{}
	blockBellatrix := consensys_eth2_openapi.BlockRequestBellatrix{}
	err := faker.FakeData(&blockBellatrix)
	t.Require().Nil(err)
	blockBellatrix.BlockHeader.Slot = "150000"
	blockBellatrix.Version = BELLATRIX

	bs.Type = "BLOCK_V2"
	bs.BeaconBlock = blockBellatrix

	body, slot, err := DecodeBeaconBlockAndSlot(ctx, bs.BeaconBlock)
	t.Require().Nil(err)
	t.Require().Equal(150000, slot)

	typeBody := body.(consensys_eth2_openapi.BlockRequestBellatrix)
	rawBytes1, err := json.Marshal(body)
	t.Require().Nil(err)

	rawBytes2, err := json.Marshal(typeBody)
	t.Require().Nil(err)
	t.Require().Equal(rawBytes1, rawBytes2)
}

func (t *WatermarkBlockDecodingTestSuite) TestDecodeCapella() {
	bs := consensys_eth2_openapi.BeaconBlockSigning{}
	blockCapella := consensys_eth2_openapi.BlockRequestCapella{
		BlockHeader: consensys_eth2_openapi.BeaconBlockHeader{},
	}
	err := faker.FakeData(&blockCapella)
	t.Require().Nil(err)
	blockCapella.BlockHeader.Slot = "200000"
	blockCapella.Version = CAPELLA

	bs.Type = "BLOCK_V2"
	bs.BeaconBlock = blockCapella

	body, slot, err := DecodeBeaconBlockAndSlot(ctx, bs.BeaconBlock)
	t.Require().Nil(err)
	t.Require().Equal(200000, slot)

	typeBody := body.(consensys_eth2_openapi.BlockRequestCapella)
	rawBytes1, err := json.Marshal(body)
	t.Require().Nil(err)

	rawBytes2, err := json.Marshal(typeBody)
	t.Require().Nil(err)
	t.Require().Equal(rawBytes1, rawBytes2)
}

func TestWatermarkBlockDecodingTestSuite(t *testing.T) {
	suite.Run(t, new(WatermarkBlockDecodingTestSuite))
}
