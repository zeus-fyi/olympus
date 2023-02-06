package athena_server

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/utils/ephemery_reset"
	"path"
	"time"
)

func StartAndConfigClientNetworkSettings(ctx context.Context, protocolNetworkID int, clientName string) {
	if protocolNetworkID == hestia_req_types.EthereumEphemeryProtocolNetworkID {
		genesisPath := dataDir.DirIn
		switch clientName {
		case client_consts.Lighthouse:
			genesisPath = path.Join(genesisPath, "/testnet")
		default:
		}

		ok, _ := ephemery_reset.Exists(path.Join(genesisPath, "/retention.vars"))
		if ok {
			kt := ephemery_reset.ExtractResetTime(path.Join(genesisPath, "/retention.vars"))
			go func(timeBeforeKill int64) {
				log.Ctx(ctx).Info().Msgf("killing ephemeral infra due to genesis reset after %d seconds", timeBeforeKill)
				time.Sleep(time.Duration(timeBeforeKill) * time.Second)
				rc := resty.New()
				// assumes you have the default choreography sidecar in your namespace cluster
				_, err := rc.R().Get("http://zeus-hydra-choreography:9999/delete/pods")
				if err != nil {
					log.Ctx(ctx).Err(err)
				}
			}(kt)
		}
	}
}
