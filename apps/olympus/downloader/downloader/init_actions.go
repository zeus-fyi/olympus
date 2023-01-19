package snapshot_init

import (
	"context"

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
		// TODO clientName is always lighthouse for validator clients for now, when you add others, add that conditional here
		err := w.DataDir.WipeDirIn()
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to wipe validator dir on startup")
			panic(err)
		}
		vsg := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{}
		vsg.ProtocolNetworkID = w.ProtocolNetworkID
		vsg.CloudCtxNs = w.CloudCtxNs
		lhW3Enable, err := artemis_validator_service_groups_models.SelectValidatorsAssignedToCloudCtxNs(ctx, vsg)
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
		err = w.DataDir.WriteToFileOutPath(ymlBytes)
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to write validators yaml")
			panic(err)
		}
	}
}
