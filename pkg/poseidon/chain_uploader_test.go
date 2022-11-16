package poseidon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/sandbox/chains"
)

type ChainUploaderTestSuite struct {
	test_suites.S3TestSuite
}

func (s *ChainUploaderTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.SetupLocalDigitalOceanS3()
	chains.ChangeToChainDataDir()
}

var brUpload = BucketRequest{
	BucketName: "zeus-fyi-ethereum",
	Protocol:   "ethereum",
	Network:    "mainnet",
	ClientType: "exec.client.standard",
	ClientName: "geth",
}

func (s *ChainUploaderTestSuite) TestChainZstdCmpAndUpload() {
	ctx := context.Background()
	pos := NewPoseidon(s.S3)
	pos.DirIn = "./ethereum/geth/data/geth"
	pos.DirOut = "./ethereum/geth_zstd_cmp"
	pos.FnIn = "geth"
	err := pos.ZstdCompressAndUpload(ctx, brUpload)
	s.Require().Nil(err)
}

func (s *ChainUploaderTestSuite) TestChainGzipCompAndUpload() {
	ctx := context.Background()
	pos := NewPoseidon(s.S3)
	pos.DirIn = "./ethereum/geth/data/geth"
	pos.DirOut = "./ethereum/geth_gzip_cmp"
	pos.FnIn = "geth"
	err := pos.GzipCompressAndUpload(ctx, brUpload)
	s.Require().Nil(err)
}

func TestChainUploaderSuite(t *testing.T) {
	suite.Run(t, new(ChainUploaderTestSuite))
}
