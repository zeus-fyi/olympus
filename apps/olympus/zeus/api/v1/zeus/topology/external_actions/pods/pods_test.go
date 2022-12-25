package pods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
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
	cliReq := ClientRequest{
		MethodHTTP:      "GET",
		Endpoint:        "health",
		Ports:           []string{"9000:9000"},
		Payload:         nil,
		EndpointHeaders: nil,
	}
	podActionRequest := PodActionRequest{
		Action:    "port-forward",
		PodName:   "eth-indexer-eth-indexer",
		ClientReq: &cliReq,
	}
	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

func (p *PodsHandlerTestSuite) TestPodPortForwardAll() {
	cliReq := ClientRequest{
		MethodHTTP:      "GET",
		Endpoint:        "health",
		Ports:           []string{"9000:9000"},
		Payload:         nil,
		EndpointHeaders: nil,
	}
	filter := string_utils.FilterOpts{DoesNotInclude: []string{"beacon", "metrics"}}
	podActionRequest := PodActionRequest{
		Action:     "port-forward-all",
		PodName:    "eth-indexer-eth-indexer",
		FilterOpts: &filter,
		ClientReq:  &cliReq,
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
	ll := zerolog.DebugLevel
	nvSize, nbSize := 100, 500
	timeout := time.Second * 30
	adminCfg := AdminConfig{
		LogLevel:                   &ll,
		ValidatorBatchSize:         &nvSize,
		ValidatorBalancesBatchSize: &nbSize,
		ValidatorBalancesTimeout:   &timeout,
	}
	cliReq := ClientRequest{
		MethodHTTP:      "POST",
		Endpoint:        "admin",
		Ports:           []string{"9000:9000"},
		Payload:         adminCfg,
		EndpointHeaders: nil,
	}
	podActionRequest := PodActionRequest{
		Action:    "port-forward",
		PodName:   "eth-indexer-eth-indexer",
		ClientReq: &cliReq,
	}

	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

func (p *PodsHandlerTestSuite) TestDescribePods() {

	kctx := zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "nyc1", Context: "do-nyc1-do-nyc1-zeus-demo", Namespace: "ephemeral"}

	tp := kns.TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: kctx,
	}
	podActionRequest := PodActionRequest{
		TopologyKubeCtxNs: tp,
		Action:            "describe",
	}
	podDescribeReq := p.postK8Request(podActionRequest, http.StatusOK, true)
	p.Require().NotEmpty(podDescribeReq.pods)
}

func (p *PodsHandlerTestSuite) TestGetPodLogs() {
	tailLines := int64(100)
	podActionRequest := PodActionRequest{
		Action:  "logs",
		PodName: "eth-indexer-eth-indexer",
		LogOpts: &v1.PodLogOptions{Container: "eth-indexer", TailLines: &tailLines},
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) TestDeletePod() {
	podActionRequest := PodActionRequest{
		Action:  "delete",
		PodName: "eth-indexer-eth-indexer",
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) TestAuditPods() {
	filter := string_utils.FilterOpts{StartsWith: "eth"}
	podActionRequest := PodActionRequest{
		Action:     "describe-audit",
		FilterOpts: &filter,
	}
	podDescribeReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	resp := podDescribeReq.logs

	var ps PodsSummary
	err := json.Unmarshal(resp, &ps)
	p.Require().NotEmpty(resp)
	p.Assert().Nil(err)
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
	p.InitLocalConfigs()
	authCfg := auth_startup.NewDefaultAuthClient(context.Background(), p.Tc.ProdLocalAuthKeysCfg)
	inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(context.Background(), authCfg)
	p.K.ConnectToK8sFromInMemFsCfgPath(inMemFs)

	z, err := p.K.GetContexts()
	p.Assert().Nil(err)
	fmt.Println(z)
	p.K.SetContext("do-nyc1-do-nyc1-zeus-demo")

	//ExternalApiPodsRoutes(eg, p.K)
}

func TestPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsHandlerTestSuite))
}
