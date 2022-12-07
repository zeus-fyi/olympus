package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacon"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestUploadWithSkeletonBaseName() {
	t.InitLocalConfigs()
	t.Eg.POST("/infra/create", CreateTopologyInfraActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	cookbooks.ChangeToCookbookDir()
	c := beacon_cookbooks.ExecClientChart
	p := beacon_cookbooks.BeaconExecClientChartPath
	c.SkeletonBaseName = "add-skeleton-base-wrw5v"
	c.Tag = "test-latest"
	resp, err := t.ZeusClient.UploadChart(ctx, p, c)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *TopologyCreateActionRequestTestSuite) TestUpload() {
	t.InitLocalConfigs()
	t.Eg.POST("/infra/create", CreateTopologyInfraActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	cookbooks.ChangeToCookbookDir()
	c := beacon_cookbooks.ExecClientChart
	p := beacon_cookbooks.BeaconExecClientChartPath

	resp, err := t.ZeusClient.UploadChart(ctx, p, c)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	c = beacon_cookbooks.ConsensusClientChart
	p = beacon_cookbooks.BeaconConsensusClientChartPath

	resp, err = t.ZeusClient.UploadChart(ctx, p, c)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
