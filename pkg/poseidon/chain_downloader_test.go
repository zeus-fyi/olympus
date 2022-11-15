package poseidon

import (
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
func (s *ChainDownloaderTestSuite) TestChainUnGzip() {
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth_gzip"
	pos.FnIn = "geth.tar.gz"
	pos.DirOut = "./ethereum/geth_ungzip"
	pos.FnOut = "data"
	err := pos.UnGzipChainData()
	s.Require().Nil(err)
}

func TestChainDownloaderTestSuite(t *testing.T) {
	suite.Run(t, new(ChainDownloaderTestSuite))
}
