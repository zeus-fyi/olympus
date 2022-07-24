package v1

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/autok8s/core"
	"github.com/zeus-fyi/olympus/pkg/client"
)

type PodsTestSuite struct {
	autok8s_core.K8TestSuite
}

func (s *PodsTestSuite) TestPodPortForward() {
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
		err := s.K.PortForwardPod(ctx, kns, "eth-indexer-eth-indexer", address, []string{ports}, startChan, stopChan)
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
	s.Require().Nil(r.Err)

	fmt.Println("end port-forwarded commands")
	fmt.Println("exiting")
}

func (s *PodsTestSuite) TestGetPods() {
	ctx := context.Background()
	var kns = autok8s_core.KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	pods, err := s.K.GetPodsUsingCtxNs(ctx, kns, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(pods)
}

func TestPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsTestSuite))
}
