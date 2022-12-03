package beacon_actions

import (
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
)

type BeaconActionsClient struct {
	zeus_client.ZeusClient
	PrintPath       filepaths.Path
	ConfigPaths     filepaths.Path
	ConsensusClient string
	ExecClient      string
}

func NewBeaconActionsClient(baseURL, bearer string) BeaconActionsClient {
	z := BeaconActionsClient{}
	z.Resty = base_rest_client.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ZeusEndpoint = "https://api.zeus.fyi"

func NewDefaultBeaconActionsClient(bearer string) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusEndpoint, bearer)
}

const ZeusLocalEndpoint = "http://localhost:9000"

func NewLocalBeaconActionsClient(bearer string) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusLocalEndpoint, bearer)
}
