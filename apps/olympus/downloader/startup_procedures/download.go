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

func ChainDownload(ctx context.Context, w WorkloadInfo) {
	log.Info().Msg("DownloadChainSnapshotRequest: Download Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)
	network := hestia_req_types.ProtocolNetworkIDToString(w.ProtocolNetworkID)
	log.Ctx(ctx).Info().Interface("network", network).Msg("DownloadChainSnapshotRequest: Downloading Chain Snapshot")
	switch w.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID, hestia_req_types.EthereumGoerliProtocolNetworkID:
		switch w.WorkloadType {
		case "beaconExecClient":
			switch w.ClientName {
			case "geth":
				log.Ctx(ctx).Info().Msg("DownloadChainSnapshotRequest: Geth Sync Starting")
				// TODO, unsure if always downloading to resync beacon is an issue or not
				b := poseidon_buckets.GethBucket(network)
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
			switch w.ClientName {
			case "lighthouse":
				log.Ctx(ctx).Info().Msg("DownloadChainSnapshotRequest: Lighthouse Sync Starting")
				b := poseidon_buckets.LighthouseBucket(network)
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
