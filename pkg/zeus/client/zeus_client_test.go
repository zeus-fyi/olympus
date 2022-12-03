package zeus_client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type ZeusClientTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient ZeusClient
}

var ctx = context.Background()

// chart workload metadata
var uploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "demo",
	ChartName:        "demo",
	ChartDescription: "demo",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var topCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "demo", // set with your own namespace
	Env:           "dev",
}

// set your own topologyID here after uploading a chart workload
var deployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669159384971627008,
	CloudCtxNs: topCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var demoChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./mocks/kubernetes_apps/demo",
	DirOut:      "./outputs/demo_read_chart",
	FnIn:        "demo", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

func (t *ZeusClientTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	//t.ZeusTestClient = NewDefaultZeusClient(tc.Bearer)
	t.ZeusTestClient = NewLocalZeusClient(tc.Bearer)
	// points working dir to inside /test
	test_base.ForceDirToTestDirLocation()

	// generates outputs to /test/outputs dir
	// uses inputs from /test/mocks dir
}

func TestZeusClientTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusClientTestSuite))
}
