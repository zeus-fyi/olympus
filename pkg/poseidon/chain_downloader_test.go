package poseidon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	"github.com/zeus-fyi/olympus/sandbox/chains"
)

type ChainDownloaderTestSuite struct {
	base.TestSuite
}

func (s *ChainDownloaderTestSuite) SetupTest() {
	s.InitLocalConfigs()
	chains.ChangeToChainDataDir()
}

var brDownload = BucketRequest{
	BucketName: "zeus-fyi",
	Protocol:   "ethereum",
	Network:    "mainnet",
	ClientType: "exec.client.standard",
	ClientName: "geth",
}

func (s *ChainDownloaderTestSuite) TestChainZstdDownloadAndDec() {
	ctx := context.Background()
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth_zstd_cmp"
	pos.FnIn = "geth.tar.zst"
	pos.DirOut = "./ethereum/geth_zstd_dec"
	pos.FnOut = "data"
	err := pos.ZstdDownloadAndDec(ctx, brDownload)
	s.Require().Nil(err)
}

func (s *ChainDownloaderTestSuite) TestChainGzipDownloadAndDec() {
	ctx := context.Background()
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth_gzip_cmp"
	pos.FnIn = "geth.tar.gz"
	pos.DirOut = "./ethereum/geth_gzip_dec"
	pos.FnOut = "data"
	err := pos.GzipDownloadAndDec(ctx, brDownload)
	s.Require().Nil(err)
}

func TestChainDownloaderTestSuite(t *testing.T) {
	suite.Run(t, new(ChainDownloaderTestSuite))
}
