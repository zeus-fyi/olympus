package olympus_snapshot_init

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func CheckForUploadBeaconJob(ctx context.Context, w WorkloadInfo) {
	ua := pg_poseidon.UploadDataDirOrchestration{
		ClientName: w.ClientName,
		OrchestrationJob: artemis_orchestrations.OrchestrationJob{
			Orchestrations: artemis_autogen_bases.Orchestrations{},
			Scheduled: artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
				Status: pg_poseidon.Pending,
			},
			CloudCtxNs: w.CloudCtxNs,
		},
	}
	shouldUpload, err := ua.CheckForPendingUploadJob(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msgf("failed to check for upload on startup", w.ClientName)
		panic(err)
	}
	if !shouldUpload {
		log.Ctx(ctx).Info().Interface("w", w).Msg("no snapshot upload job found, skipping")
		return
	}
	log.Ctx(ctx).Info().Interface("w", w).Msg("disk upload job found, starting snapshot upload")
	UploadSelector(ctx, w)
	err = ua.MarkUploadComplete(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msgf("failed to mark upload complete for %s", w.ClientName)
		panic(err)
	}
}

func UploadSelector(ctx context.Context, w WorkloadInfo) {
	network := hestia_req_types.ProtocolNetworkIDToString(w.ProtocolNetworkID)
	log.Info().Interface("network", network).Msg("UploadChainSnapshotRequest: Upload Sync Starting")
	pos := poseidon.NewPoseidon(athena.AthenaS3Manager)

	switch w.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID, hestia_req_types.EthereumGoerliProtocolNetworkID:
		switch w.WorkloadType {
		case "beaconExecClient":
			switch w.ClientName {
			case "geth":
				log.Ctx(ctx).Info().Msg("UploadChainSnapshotRequest: Geth Upload Starting")
				b := poseidon_buckets.GethBucket(network)
				pos.FnIn = b.GetBucketKey()
				err := pos.Lz4CompressAndUpload(ctx, b)
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
				pos.FnIn = b.GetBucketKey()
				err := pos.Lz4CompressAndUpload(ctx, b)
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
		panic(err)
	}
}
