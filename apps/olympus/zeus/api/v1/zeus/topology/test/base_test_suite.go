package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	v1 "k8s.io/api/core/v1"
)

var Kns autok8s_core.KubeCtxNs

type TopologyActionRequestTestSuite struct {
	E *echo.Echo
	autok8s_core.K8TestSuite
	DB test_suites.PGTestSuite
	Ts chronos.Chronos
}

type TestResponse struct {
	logs []byte
	pods v1.PodList
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
