package olympus_snapshot_init

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func ChainDownload(ctx context.Context) {
	log.Info().Msg("DownloadChainSnapshotRequest: Download Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)

	switch Workload.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		switch Workload.WorkloadType {
		case "beaconExecClient":
			switch Workload.ClientName {
			case "geth":
				// TODO, unsure if always downloading to resync beacon is an issue or not
				b := poseidon_buckets.GethMainnetBucket
				err := pos.SyncDownload(ctx, b)
				if err != nil {
					log.Ctx(ctx).Err(err)
					panic(err)
				}
			default:
				err := errors.New("invalid client workload type")
				log.Ctx(ctx).Err(err)
			}
		case "beaconConsensusClient":
			switch Workload.ClientName {
			case "lighthouse":
				b := poseidon_buckets.LighthouseMainnetBucket
				err := pos.SyncDownload(ctx, b)
				if err != nil {
					log.Ctx(ctx).Err(err)
					panic(err)
				}
			default:
				err := errors.New("invalid client workload type")
				log.Ctx(ctx).Err(err)
			}
		default:
			err := errors.New("invalid client workload type")
			log.Ctx(ctx).Err(err)
		}
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
	default:
		err := errors.New("invalid or unsupported protocol network id")
		log.Ctx(ctx).Err(err)
	}
}
