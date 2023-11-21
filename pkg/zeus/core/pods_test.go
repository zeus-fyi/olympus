package zeus_core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type PodsTestSuite struct {
	K8TestSuite
}

func (s *PodsTestSuite) TestPodPortForward() {
	c := client.Client{}
	c.E = "http://localhost:9000"

	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "zeus-k8s-blockchain", Namespace: "eth-indexer"}

	address := "localhost"
	ports := "9000:9000"

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	go func() {
		fmt.Println("start port-forward thread")
		err := s.K.PortForwardPod(ctx, kns, "eth-indexer-eth-indexer", address, []string{ports}, startChan, stopChan, nil)
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

	//fmt.Println("do port-forwarded commands")
	//r := c.Get(ctx, "http://localhost:9000/health")
	//s.Require().Nil(r.Err)
	//
	//fmt.Println("end port-forwarded commands")
	//fmt.Println("exiting")
}

func (s *PodsTestSuite) TestDeletePods() {
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral"}
	filter := strings_filter.FilterOpts{
		DoesNotStartWithThese: nil,
		StartsWithAnyOfThese:  nil,
		StartsWith:            "",
		Contains:              "client",
		DoesNotInclude:        nil,
	}
	err := s.K.DeleteAllPodsLike(ctx, kns, "", nil, &filter)
	s.Require().Nil(err)
}

func (s *PodsTestSuite) TestGetPods() {
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral"}

	pods, err := s.K.GetPodsUsingCtxNs(ctx, kns, nil, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(pods)

	for _, pod := range pods.Items {
		fmt.Println(pod.GetName())
	}
}

func TestPodsTestSuite(t *testing.T) {
	suite.Run(t, new(PodsTestSuite))
}
