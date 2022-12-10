package zeus_config_map_reqs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

const (
	KeySwapAction              = "key-swap"
	SetOrCreateKeyFromExisting = "set-or-create-from-key"
)

type ConfigMapActionRequest struct {
	zeus_req_types.TopologyDeployRequest
	Action        string
	ConfigMapName string
	Keys          KeySwap
	FilterOpts    *string_utils.FilterOpts
}

// KeySwap If using create new key from existing then keyOne=keyToCopy, keyTwo=keyToSetOrCreateFromCopy
type KeySwap struct {
	KeyOne string `json:"keyOne"`
	KeyTwo string `json:"keyTwo"`
}
