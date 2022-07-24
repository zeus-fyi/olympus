package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/autok8s/core"
	"github.com/zeus-fyi/olympus/pkg/client"
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

func (p *PodsHandlerTestSuite) TestPodPortForward() {
	c := client.Client{}
	c.E = "http://localhost:9000"

	ctx := context.Background()
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	address := "localhost"
	ports := "9000:9000"

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	go func() {
		fmt.Println("start port-forward thread")
		err := p.K.PortForwardPod(ctx, kns, "eth-indexer-eth-indexer", address, []string{ports}, startChan, stopChan)
		fmt.Println(err)
		fmt.Println("done port-forward")
	}()

	fmt.Println("awaiting signal")
	<-startChan
	defer close(stopChan)
	fmt.Println("port ready chan ok")
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		close(stopChan)
	}()

	fmt.Println("do port-forwarded commands")
	r := c.Get(ctx, "http://localhost:9000/health")
	p.Require().Nil(r.Err)

	fmt.Println("end port-forwarded commands")
	fmt.Println("exiting")
}

func (p *PodsHandlerTestSuite) TestGetPods() {
	ctx := context.Background()
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	pods, err := p.K.GetPodsUsingCtxNs(ctx, kns, nil)
	p.Require().Nil(err)
	p.Require().NotEmpty(pods)
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
	p.E = InitRouter(e, p.K)
}

func TestPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsHandlerTestSuite))
}
