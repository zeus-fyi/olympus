package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	v1 "k8s.io/api/core/v1"
)

var Kns autok8s_core.KubeCtxNs

var TestOrgUser = org_users.OrgUser{autogen_bases.OrgUsers{
	OrgID:  1667266332674446258,
	UserID: 1667266332670878528,
}}

var TestTopologyID = 6951056435719556916

type TopologyActionRequestTestSuite struct {
	E *echo.Echo
	autok8s_core.K8TestSuite
	D  test_suites.DatastoresTestSuite
	Ts chronos.Chronos
}

type TestResponse struct {
	logs []byte
	pods v1.PodList
}

func (t *TopologyActionRequestTestSuite) SetupTest() {
	t.ConnectToK8s()
	t.InitLocalConfigs()

	t.D.PGTest.SetupPGConn()
	t.D.PG = t.D.PGTest.Pg
}

func (t *TopologyActionRequestTestSuite) PostTopologyRequest(topologyActionRequest base.TopologyActionRequest, httpCode int) TestResponse {
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
