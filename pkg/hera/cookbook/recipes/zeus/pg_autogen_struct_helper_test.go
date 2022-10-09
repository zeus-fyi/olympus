package zeus

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook"
)

type ZeusRecipeTestSuite struct {
	cookbook.CookbookTestSuiteBase
}

func (s *ZeusRecipeTestSuite) TestZeusDerivativeStructGen() {
	baseFw.PackageName = "zeus.go"
	err := GenerateZeusStruct(baseFw)
	s.Require().Nil(err)
}

func TestZeusRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusRecipeTestSuite))
}
