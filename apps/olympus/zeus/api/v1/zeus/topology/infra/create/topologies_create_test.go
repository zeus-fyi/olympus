package create_infra

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacon"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
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

	c := beacon_cookbooks.ConsensusClientChart

	createRequest := TopologyCreateRequest{
		TopologyName:     c.TopologyName,
		ChartName:        c.ChartName,
		ChartDescription: c.ChartDescription,
		Version:          fmt.Sprintf("v0.0.%d", +t.Ts.UnixTimeStampNow()),

		SkeletonBaseID: 0,
	}
	cookbooks.ChangeToCookbookDir()

	// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
	p := beacon_cookbooks.BeaconConsensusClientChartPath
	comp := compression.NewCompression()
	err := comp.GzipCompressDir(&p)
	t.Require().Nil(err)

	resp, err := t.ZeusClient.UploadChart(ctx, p, zeus_req_types.TopologyCreateRequest(createRequest))
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
