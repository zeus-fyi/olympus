package echo_api_app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func (s *MainTestSuite) TestMainCodeGen() {
	resp := genFile()
	s.Assert().NotEmpty(resp)
	fmt.Printf("%#v", resp)

	err := resp.Save("m.go")
	s.Assert().Nil(err)
}
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
