package olympus_beacon

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

func (t *ZeusAppsTestSuite) TestRead() {
	_, err := t.ZeusTestClient.ReadTopologies(ctx)
	t.Require().Nil(err)
}

// chart workload metadata
func newUploadChart(name string) zeus_req_types.TopologyCreateRequest {
	return zeus_req_types.TopologyCreateRequest{
		TopologyName:     name,
		ChartName:        name,
		ChartDescription: name,
		Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
}
