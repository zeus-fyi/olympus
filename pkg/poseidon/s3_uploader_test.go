package poseidon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type ChainUploaderTestSuite struct {
	test_suites_s3.S3TestSuite
}

var brUpload = BucketRequest{
	BucketName: "zeusfyi",
}

func (s *ChainUploaderTestSuite) TestOvHTextFileZstdCmpAndUpload() {
	ctx := context.Background()
	pos := NewS3Poseidon(s.OvhS3)
	pos.DirIn = "./"
	pos.FnIn = "tmp.txt"
	err := pos.ZstCompressFile(ctx, &pos.Path)
	s.Require().Nil(err)
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
