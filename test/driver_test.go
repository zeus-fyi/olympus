package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type DriverTestSuite struct {
	test_suites.TemporalTestSuite
}

func (t *DriverTestSuite) TestDrive() {

	//start := make(chan struct{}, 1)
	//go func() {
	//	close(start)
	//	zeus_server.Zeus()
	//}()
	//
	//<-start

	err := CallAPI()
	t.Require().Nil(err)
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
