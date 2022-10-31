package infra

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	v12 "github.com/zeus-fyi/olympus/zeus/api/v1"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyActionRequestTestSuite) SetupTest() {
	e := echo.New()
	t.K.CfgPath = t.K.DefaultK8sCfgPath()
	t.K.ConnectToK8s()
	t.DB.SetupPGConn()
	t.E = v12.InitRouter(e, t.K)
}

func TestTopologyActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyActionRequestTestSuite))
}
