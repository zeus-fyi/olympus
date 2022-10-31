package topology

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	v12 "github.com/zeus-fyi/olympus/zeus/api/v1"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	v1 "k8s.io/api/core/v1"
)

type TestResponse struct {
	logs []byte
	pods v1.PodList
}

var kns autok8s_core.KubeCtxNs

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

//func (t *TopologyActionRequestTestSuite) TestReadChart() {
//	topologyActionRequest := base.TopologyActionRequest{
//		Action:     "read",
//		K8sRequest: zeus_pkg.K8sRequest{Kns: kns},
//		Cluster:    clusters.NewCluster(),
//	}
//	t.postTopologyRequest(topologyActionRequest, 200)
//}

func TestTopologyActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyActionRequestTestSuite))
}
