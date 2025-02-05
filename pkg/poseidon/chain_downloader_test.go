package poseidon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type ChainDownloaderTestSuite struct {
	test_suites_s3.S3TestSuite
}

func (s *ChainDownloaderTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.SetupLocalDigitalOceanS3()
}

var brDownload = BucketRequest{
	BucketName: "zeus-fyi-ethereum",
	Protocol:   "ethereum",
	Network:    "mainnet",
	ClientType: "exec.client.standard",
	ClientName: "geth",
}

func (s *ChainDownloaderTestSuite) TestChainZstdDownloadAndDec() {
	ctx := context.Background()
	pos := NewPoseidon(s.S3)
	pos.DirIn = "./ethereum/geth_zstd_download"
	pos.FnIn = "geth.tar.zst"
	pos.DirOut = "./ethereum/geth_zstd_download"
	pos.FnOut = "geth"
	err := pos.ZstdDownloadAndDec(ctx, brDownload)
	s.Require().Nil(err)
}

func (s *ChainDownloaderTestSuite) TestChainGzipDownloadAndDec() {
	ctx := context.Background()
	pos := NewPoseidon(s.S3)
	pos.DirIn = "./ethereum/geth_gzip_download"
	pos.FnIn = "geth.tar.gz"
	pos.DirOut = "./ethereum/geth_gzip_dec"
	pos.FnOut = "geth"
	err := pos.GzipDownloadAndDec(ctx, brDownload)
	s.Require().Nil(err)
}

func TestChainDownloaderTestSuite(t *testing.T) {
	suite.Run(t, new(ChainDownloaderTestSuite))
}
