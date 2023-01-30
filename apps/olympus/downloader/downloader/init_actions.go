package snapshot_init

import (
	"context"
	"path"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/zeus/pkg/utils/ephemery_reset"
)

func InitWorkloadAction(ctx context.Context, w WorkloadInfo) {
	// the below uses a switch case to download if an ephemeralClientName is used
	ephemery_reset.ExtractAndDecEphemeralTestnetConfig(Workload.DataDir, w.ClientName)
	switch w.WorkloadType {
	case "validatorClient":
		log.Ctx(ctx).Info().Msg("starting validators sync")
		// TODO clientName is always lighthouse for validator clients for now, when you add others, add that conditional here
		err := w.DataDir.WipeDirIn()
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to wipe validator dir on startup")
			panic(err)
		}
		vsg := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{}
		vsg.ProtocolNetworkID = w.ProtocolNetworkID
		vsg.ValidatorClientNumber = w.ReplicaCountNum
		lhW3Enable, err := artemis_validator_service_groups_models.SelectValidatorsAssignedToCloudCtxNs(ctx, vsg, w.CloudCtxNs)
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to select validators")
			panic(err)
		}
		ymlBytes, err := yaml.Marshal(&lhW3Enable)
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to marshall yaml")
			panic(err)
		}
		w.DataDir.FnOut = "validator_definitions.yaml"
		w.DataDir.DirOut = path.Join(w.DataDir.DirIn, "/validators")
		err = w.DataDir.WriteToFileOutPath(ymlBytes)
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to write validators yaml")
			panic(err)
		}
		log.Ctx(ctx).Info().Msg("validators sync complete")
	case "beaconExecClient", "beaconConsensusClient":
		log.Ctx(ctx).Info().Msg("starting chain sync")
		ChainDownload(ctx)
		log.Ctx(ctx).Info().Msg("chain sync complete")
	}
}
