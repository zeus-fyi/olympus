package template_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestCodeGen() {

}
func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(FuncTestSuite))
}
