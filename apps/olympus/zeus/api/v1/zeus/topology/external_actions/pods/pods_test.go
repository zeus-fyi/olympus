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
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
	v1 "k8s.io/api/core/v1"
)

var ctx = context.Background()

type PodsHandlerTestSuite struct {
	test.TopologyActionRequestTestSuite
}

type TestResponse struct {
	logs []byte
	pods v1.PodList
}

func (p *PodsHandlerTestSuite) TestPodPortForwardGET() {
	cliReq := zeus_pods_reqs.ClientRequest{
		MethodHTTP:      "GET",
		Endpoint:        "health",
		Ports:           []string{"9000:9000"},
		Payload:         nil,
		EndpointHeaders: nil,
	}
	podActionRequest := zeus_pods_reqs.PodActionRequest{
		Action:    "port-forward",
		PodName:   "eth-indexer-eth-indexer",
		ClientReq: &cliReq,
	}
	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

func (p *PodsHandlerTestSuite) TestPodPortForwardAll() {
	cliReq := zeus_pods_reqs.ClientRequest{
		MethodHTTP:      "GET",
		Endpoint:        "health",
		Ports:           []string{"9000:9000"},
		Payload:         nil,
		EndpointHeaders: nil,
	}
	filter := strings_filter.FilterOpts{DoesNotInclude: []string{"beacon", "metrics"}}
	podActionRequest := zeus_pods_reqs.PodActionRequest{
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
	cliReq := zeus_pods_reqs.ClientRequest{
		MethodHTTP:      "POST",
		Endpoint:        "admin",
		Ports:           []string{"9000:9000"},
		Payload:         adminCfg,
		EndpointHeaders: nil,
	}
	podActionRequest := zeus_pods_reqs.PodActionRequest{
		Action:    "port-forward",
		PodName:   "eth-indexer-eth-indexer",
		ClientReq: &cliReq,
	}

	podPortForwardReq := p.postK8Request(podActionRequest, http.StatusOK, false)
	p.Require().NotEmpty(podPortForwardReq.logs)
}

//
//func (p *PodsHandlerTestSuite) TestDescribePods() {
//	p.InitLocalConfigs()
//	p.Eg.POST("/pods", HandlePodActionRequest)
//	start := make(chan struct{}, 1)
//	go func() {
//		close(start)
//		_ = p.E.Start(":9010")
//	}()
//
//	<-start
//	defer p.E.Shutdown(ctx)
//	kctx := zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral-staking"}
//
//	podActionRequest := zeus_pods_reqs.PodActionRequest{
//		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
//			TopologyID: 0,
//			CloudCtxNs: kctx,
//		},
//		Action: "describe",
//	}
//	resp, err := p.ZeusClient.GetPods(ctx, podActionRequest)
//	p.Require().NoError(err)
//	p.Require().NotEmpty(resp)
//}

func (p *PodsHandlerTestSuite) TestGetPodLogs() {
	tailLines := int64(100)
	podActionRequest := zeus_pods_reqs.PodActionRequest{
		Action:  "logs",
		PodName: "eth-indexer-eth-indexer",
		LogOpts: &v1.PodLogOptions{Container: "eth-indexer", TailLines: &tailLines},
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) TestDeletePod() {
	podActionRequest := zeus_pods_reqs.PodActionRequest{
		Action:  "delete",
		PodName: "eth-indexer-eth-indexer",
	}
	p.postK8Request(podActionRequest, http.StatusOK, false)
}

func (p *PodsHandlerTestSuite) TestAuditPods() {
	p.InitLocalConfigs()
	p.Eg.POST("/pods", HandlePodActionRequest)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = p.E.Start(":9010")
	}()

	<-start
	defer p.E.Shutdown(ctx)
	kctx := zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral-staking"}

	podActionRequest := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			TopologyID: 0,
			CloudCtxNs: kctx,
		},
		Action: "describe-audit"}
	topologyActionRequestPayload, err := json.Marshal(podActionRequest)
	p.Require().NoError(err)

	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))
	//resp, err := p.ZeusClient.GetPodsAudit(ctx, podActionRequest)
	//p.Require().NoError(err)
	//p.Require().NotEmpty(resp)
}

func (p *PodsHandlerTestSuite) postK8Request(podActionRequest zeus_pods_reqs.PodActionRequest, httpCode int, unmarshall bool) TestResponse {
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
