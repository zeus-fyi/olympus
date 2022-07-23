package autok8s_core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/client"
)

type PodsTestSuite struct {
	K8TestSuite
}

func (s *PodsTestSuite) TestPodPortForward() {
	c := client.Client{}

	ctx := context.Background()
	var kns = KubeCtxNs{Env: "", CloudProvider: "", Region: "", CtxType: "data", Namespace: "eth-indexer"}

	address := "localhost"
	ports := "8080:8080"

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	go func() {
		fmt.Println("start port-forward thread")
		err := s.k.PortForwardPod(ctx, kns, "ethereum-qt-primary-beacon", address, []string{ports}, startChan, stopChan)
		fmt.Println(err)
		fmt.Println("done port-forward")
	}()

	fmt.Println("awaiting signal")
	<-startChan
	fmt.Println("port ready chan ok")
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		close(stopChan)
	}()

	fmt.Println("do port-forwarded commands")
	r := c.Get(ctx, "metrics")
	s.Require().Nil(r.Err)
	fmt.Println("end port-forwarded commands")
	fmt.Println("exiting")
}

func (s *PodsTestSuite) TestGetPods() {
	ctx := context.Background()
	var kns = KubeCtxNs{Env: "", CloudProvider: "", Region: "", CtxType: "data", Namespace: "eth-indexer"}

	pods, err := s.k.GetPodsUsingCtxNs(ctx, kns, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(pods)
}

func TestKPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsTestSuite))
}
