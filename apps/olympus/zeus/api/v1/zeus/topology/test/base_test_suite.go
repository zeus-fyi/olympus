package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/pretty"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
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
	D        test_suites.DatastoresTestSuite
	Ts       chronos.Chronos
	Endpoint string
}

type TestResponse struct {
	Logs []byte
	pods v1.PodList
}

func (t *TopologyActionRequestTestSuite) SetupTest() {
	t.ConnectToK8s()
	t.InitLocalConfigs()

	t.D.PGTest.SetupPGConn()
	t.D.PG = t.D.PGTest.Pg
	t.E = echo.New()
}

func (t *TopologyActionRequestTestSuite) AddEndpointHandler(h echo.HandlerFunc) {
	t.E.POST(t.Endpoint, h)
}

func (t *TopologyActionRequestTestSuite) PostTopologyRequest(topologyActionRequest interface{}, httpCode int) TestResponse {
	topologyActionRequestPayload, err := json.Marshal(topologyActionRequest)
	t.Assert().Nil(err)

	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))
	req := httptest.NewRequest(http.MethodPost, t.Endpoint, strings.NewReader(string(topologyActionRequestPayload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	t.E.ServeHTTP(rec, req)
	t.Equal(httpCode, rec.Code)

	var tr TestResponse
	tr.Logs = rec.Body.Bytes()
	fmt.Println("resp json")
	t.Assert().Nil(err)

	result := pretty.Pretty(rec.Body.Bytes())
	result = pretty.Color(result, pretty.TerminalStyle)

	fmt.Println(string(result))

	return tr
}
