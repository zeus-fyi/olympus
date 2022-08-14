package _func

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestMainCodeGen() {
	resp := genFile()
	s.Assert().NotEmpty(resp)
	fmt.Printf("%#v", resp)

	err := resp.Save("function.go")
	s.Assert().Nil(err)
}
func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(FuncTestSuite))
}
