package zeus_config_map_reqs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

const (
	KeySwapAction = "key-swap"
)

type ConfigMapActionRequest struct {
	zeus_req_types.TopologyDeployRequest
	Action        string
	ConfigMapName string
	Keys          KeySwap
	FilterOpts    *string_utils.FilterOpts
}

type KeySwap struct {
	KeyOne string `json:"keyOne"`
	KeyTwo string `json:"keyTwo"`
}
