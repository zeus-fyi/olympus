package cli_wrapper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CliWrapperTestSuite struct {
	suite.Suite
}

func (s *CliWrapperTestSuite) SetupTest() {
}

func (s *CliWrapperTestSuite) TestEcho() {
	hw := "Hello World"
	cmd := TaskCmd{Command: "echo", Args: []string{hw}}
	stdOut, stdError, err := cmd.ExecuteCmd()
	s.Equal(fmt.Sprintf("%s\n", hw), stdOut)
	s.Nil(err)
	s.Empty(stdError)
}

func TestCliWrapperTestSuite(t *testing.T) {
	suite.Run(t, new(CliWrapperTestSuite))
}
