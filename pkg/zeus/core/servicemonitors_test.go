package zeus_core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/servicemonitors"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ServiceMonitorsTestSuite struct {
	K8TestSuite
}

func (s *ServiceMonitorsTestSuite) TestGetServiceMonitor() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral"}

	sm, err := s.K.GetServiceMonitor(ctx, kns, "zeus-consensus-client-monitor", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(sm)
}

func (s *ServiceMonitorsTestSuite) TestCreateServiceMonitor() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral"}

	sm := servicemonitors.NewServiceMonitor()
	filepath := "/Users/alex/go/Olympus/olympus/cookbooks/olympus/ethereum/beacons/infra/consensus_client/servicemonitor.yaml"
	yr := transformations.YamlFileIO{}
	jsonBytes, err := yr.ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &sm.K8sServiceMonitor)
	s.Require().Nil(err)
	s.Require().NotEmpty(sm.K8sServiceMonitor)
	svc, err := s.K.CreateServiceMonitor(ctx, kns, &sm.K8sServiceMonitor, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)
}

func TestServiceMonitorsTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceMonitorsTestSuite))
}
