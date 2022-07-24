package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/autok8s/core"
	v1 "k8s.io/api/core/v1"
)

type PodsHandlerTestSuite struct {
	E *echo.Echo
	autok8s_core.K8TestSuite
}

type TestResponse struct {
	logs []byte
	pods v1.PodList
}

func (p *PodsHandlerTestSuite) SetupTest() {
	p.SetupTestServer()
}

func (p *PodsHandlerTestSuite) TestPodPortForwardGET() {
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	cliReq := ClientRequest{
		MethodHTTP:      "GET",
		Endpoint:        "health",
		Ports:           []string{"9000:9000"},
		Payload:         nil,
		EndpointHeaders: nil,
	}
	podActionRequest := PodActionRequest{
		Action:     "port-forward",
		PodName:    "eth-indexer-eth-indexer",
		ClientReq:  &cliReq,
		K8sRequest: K8sRequest{kns},
	}
	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

type AdminConfig struct {
	LogLevel *zerolog.Level

	ValidatorBatchSize         *int
	ValidatorBalancesBatchSize *int
	ValidatorBalancesTimeout   *time.Duration
}

func (p *PodsHandlerTestSuite) TestPodPortForwardPOST() {
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	ll := zerolog.DebugLevel
	nvSize, nbSize := 100, 500
	timeout := time.Second * 30
	adminCfg := AdminConfig{
		LogLevel:                   &ll,
		ValidatorBatchSize:         &nvSize,
		ValidatorBalancesBatchSize: &nbSize,
		ValidatorBalancesTimeout:   &timeout,
	}
	payload, err := json.Marshal(adminCfg)
	p.Require().Nil(err)
	payloadStr := string(payload)
	cliReq := ClientRequest{
		MethodHTTP:      "POST",
		Endpoint:        "admin",
		Ports:           []string{"9000:9000"},
		Payload:         &payloadStr,
		EndpointHeaders: nil,
	}
	podActionRequest := PodActionRequest{
		Action:     "port-forward",
		PodName:    "eth-indexer-eth-indexer",
		ClientReq:  &cliReq,
		K8sRequest: K8sRequest{kns},
	}

	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

func (p *PodsHandlerTestSuite) TestDescribePods() {
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	podActionRequest := PodActionRequest{
		Action:     "describe",
		K8sRequest: K8sRequest{kns},
	}
	podDescribeReq := p.postK8Request(podActionRequest, http.StatusOK, true)
	p.Require().NotEmpty(podDescribeReq.pods)
}

func (p *PodsHandlerTestSuite) TestGetPodLogs() {
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	tailLines := int64(100)
	podActionRequest := PodActionRequest{
		Action:     "logs",
		PodName:    "eth-indexer-eth-indexer",
		K8sRequest: K8sRequest{kns},
		LogOpts:    &v1.PodLogOptions{Container: "eth-indexer", TailLines: &tailLines},
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) TestDeletePod() {
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	podActionRequest := PodActionRequest{
		Action:     "delete",
		PodName:    "eth-indexer-eth-indexer",
		K8sRequest: K8sRequest{kns},
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) postK8Request(podActionRequest PodActionRequest, httpCode int, unmarshall bool) TestResponse {
	podActionRequestPayload, err := json.Marshal(podActionRequest)
	p.Assert().Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/pods", strings.NewReader(string(podActionRequestPayload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	p.E.ServeHTTP(rec, req)
	p.Equal(httpCode, rec.Code)

	var tr TestResponse

	if unmarshall {
		err = json.Unmarshal(rec.Body.Bytes(), &tr.pods)
		p.Require().Nil(err)
		return tr
	} else {
		tr.logs = rec.Body.Bytes()
	}
	return tr
}

func (p *PodsHandlerTestSuite) SetupTestServer() {
	e := echo.New()
	p.K.CfgPath = p.K.DefaultK8sCfgPath()
	p.K.ConnectToK8s()
	p.E = InitRouter(e, p.K)
}

func TestPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsHandlerTestSuite))
}
