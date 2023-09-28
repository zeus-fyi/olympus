package hestia_digitalocean

import (
	"fmt"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DoKubernetesTestSuite struct {
	test_suites_base.TestSuite
	do DigitalOcean
}

func (s *DoKubernetesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.do = InitDoClient(ctx, s.Tc.DigitalOceanAPIKey)
	s.Require().NotNil(s.do.Client)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
func (s *DoKubernetesTestSuite) TestGetClusterContexts() {
	k8sContext, _, err := s.do.Client.Kubernetes.List(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotEmpty(k8sContext)
}
func (s *DoKubernetesTestSuite) TestGetNodePools() {
	nycContext := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"
	nodePools, _, err := s.do.Client.Kubernetes.ListNodePools(ctx, nycContext, nil)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodePools)

	for _, np := range nodePools {
		fmt.Println(np.ID)
	}
}

func (s *DoKubernetesTestSuite) TestCreateNodePool() {
	clusterUUID := uuid.New()

	taint := godo.Taint{
		Key:    fmt.Sprintf("org-%d", s.Tc.ProductionLocalTemporalOrgID),
		Value:  fmt.Sprintf("org-%d", s.Tc.ProductionLocalTemporalOrgID),
		Effect: "NoSchedule",
	}

	var labels map[string]string
	labels = make(map[string]string)
	labels = AddDoNvmeLabels(labels)
	suffix := strings.Split(clusterUUID.String(), "-")[0]
	nodesReq := &godo.KubernetesNodePoolCreateRequest{
		Name:   fmt.Sprintf("nodepool-%d-%s", s.Tc.ProductionLocalTemporalOrgID, suffix),
		Size:   "so1_5-4vcpu-32gb",
		Count:  int(1),
		Tags:   nil,
		Labels: labels,
		Taints: []godo.Taint{taint},
	}

	clusterID := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"
	np, err := s.do.CreateNodePool(ctx, clusterID, nodesReq)
	s.Require().NoError(err)
	s.Require().NotNil(np)
}

func (s *DoKubernetesTestSuite) TestDeleteNodePool() {
	//// TODO
	//err = s.do.RemoveNodePool(ctx, clusterID, np.ID)
	//s.Require().NoError(err)
}

func TestDoKubernetesTestSuite(t *testing.T) {
	suite.Run(t, new(DoKubernetesTestSuite))
}
