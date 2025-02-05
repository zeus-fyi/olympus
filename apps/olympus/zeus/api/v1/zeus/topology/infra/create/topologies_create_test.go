package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	olympus_hydra_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra"
	hephaestus_olympus_cookbook "github.com/zeus-fyi/olympus/cookbooks/olympus/hephaestus"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	beacon_cookbooks "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestInternalChartUploadJobs() {
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

	olympus_cookbooks.ChangeToCookbookDir()
	cdCfg := hephaestus_olympus_cookbook.HephaestusClusterDefinition

	bcs, err := cdCfg.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)

	for _, bc := range bcs {
		resp, uerr := t.ZeusClient.UploadChart(ctx, bc.SkeletonBaseNameChartPath, bc.SkeletonBaseChart)
		t.Require().Nil(uerr)
		t.Assert().NotEmpty(resp)
	}
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

	olympus_cookbooks.ChangeToCookbookDir()
	c := beacon_cookbooks.ExecClientChart
	p := beacon_cookbooks.BeaconExecClientChartPath
	c.ClusterClassName = "lz2l2xd6wk"
	c.ComponentBaseName = "test-cluster-base"
	c.SkeletonBaseName = "whatever"
	c.Tag = "latest"
	resp, err := t.ZeusClient.UploadChart(ctx, p, c)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *TopologyCreateActionRequestTestSuite) TestInternalChartUpload() {
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

	olympus_cookbooks.ChangeToCookbookDir()
	cdCfg := olympus_hydra_cookbooks.HydraClusterConfig(&olympus_hydra_cookbooks.HydraClusterDefinition, "ephemery")

	bcs, err := cdCfg.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)

	for _, bc := range bcs {
		resp, uerr := t.ZeusClient.UploadChart(ctx, bc.SkeletonBaseNameChartPath, bc.SkeletonBaseChart)
		t.Require().Nil(uerr)
		t.Assert().NotEmpty(resp)
	}
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

	olympus_cookbooks.ChangeToCookbookDir()
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
