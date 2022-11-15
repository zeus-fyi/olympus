package poseidon

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	"github.com/zeus-fyi/olympus/sandbox/chains"
)

type ChainUploaderTestSuite struct {
	base.TestSuite
}

func (s *ChainUploaderTestSuite) SetupTest() {
	s.InitLocalConfigs()
	chains.ChangeToChainDataDir()
}

func (s *ChainUploaderTestSuite) TestChainZstdComp() {
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth/data/geth"
	pos.DirOut = "./ethereum/geth_zstd_cmp"
	pos.FnIn = "geth"
	err := pos.ZstdCompressChainData()
	s.Require().Nil(err)
}

func (s *ChainUploaderTestSuite) TestChainGzipComp() {
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth/data/geth"
	pos.DirOut = "./ethereum/geth_gzip_cmp"
	pos.FnIn = "geth"
	err := pos.GzipCompressChainData()
	s.Require().Nil(err)
}

func TestChainUploaderSuite(t *testing.T) {
	suite.Run(t, new(ChainUploaderTestSuite))
}
