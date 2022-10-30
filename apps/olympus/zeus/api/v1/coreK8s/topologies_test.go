package coreK8s

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type TopologyActionRequestTestSuite struct {
	E *echo.Echo
	autok8s_core.K8TestSuite
	DB test_suites.PGTestSuite
}

func (t *TopologyActionRequestTestSuite) SetupTest() {
	e := echo.New()
	t.K.CfgPath = t.K.DefaultK8sCfgPath()
	t.K.ConnectToK8s()
	t.DB.SetupPGConn()
	t.E = InitRouter(e, t.K)
}

func (t *TopologyActionRequestTestSuite) TestChartQueryHandler() {
	topologyActionRequest := TopologyActionRequest{
		Action:     "read",
		K8sRequest: K8sRequest{Kns: kns},
		Cluster:    clusters.NewCluster(),
	}

	t.postTopologyRequest(topologyActionRequest, 200)
}

func (t *TopologyActionRequestTestSuite) postTopologyRequest(topologyActionRequest TopologyActionRequest, httpCode int) TestResponse {
	topologyActionRequestPayload, err := json.Marshal(topologyActionRequest)
	t.Assert().Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/topology", strings.NewReader(string(topologyActionRequestPayload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	t.E.ServeHTTP(rec, req)
	t.Equal(httpCode, rec.Code)

	var tr TestResponse
	tr.logs = rec.Body.Bytes()
	return tr
}
func TestTopologyActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyActionRequestTestSuite))
}
