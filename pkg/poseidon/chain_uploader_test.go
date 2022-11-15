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
func (s *ChainUploaderTestSuite) TestChainGzip() {
	pos := NewPoseidon()
	pos.DirIn = "./ethereum/geth/data/geth"
	pos.DirOut = "./ethereum/geth_gzip"
	pos.FnIn = "geth"
	err := pos.GzipChainData()
	s.Require().Nil(err)
}

func TestChainUploaderSuite(t *testing.T) {
	suite.Run(t, new(ChainUploaderTestSuite))
}
