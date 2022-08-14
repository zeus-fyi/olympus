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
}
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
