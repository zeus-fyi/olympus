package olympus_snapshot_init

import (
	"context"
	init_jwt "github.com/zeus-fyi/zeus/pkg/aegis/jwt"
	"path"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	beacon_cookbooks "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/utils/ephemery_reset"
)

func InitWorkloadAction(ctx context.Context, w WorkloadInfo) {
	log.Ctx(ctx).Info().Interface("w", w).Msg("init workload action")
	switch w.WorkloadType {
	case "validatorClient":
		log.Ctx(ctx).Info().Msg("starting validators sync")
		// TODO clientName is always lighthouse for validator clients for now, when you add others, add that conditional here
		err := w.DataDir.WipeDirIn()
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to wipe validator dir on startup")
			panic(err)
		}
		EphemeryReset()
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
		EphemeryReset()
		log.Ctx(ctx).Info().Msg("starting chain sync")
		ChainDownload(ctx)
		log.Ctx(ctx).Info().Msg("chain sync complete")
		if useDefaultToken {
			log.Ctx(ctx).Info().Msg("setting jwt token to default")
			err := init_jwt.SetTokenToDefault(w.DataDir, "jwt.hex", jwtToken)
			if err != nil {
				log.Ctx(ctx).Panic().Err(err).Msg("failed to set jwt token to default")
				panic(err)
			}
		}
	}
}

func EphemeryReset() {
	if Workload.ProtocolNetworkID == hestia_req_types.EthereumEphemeryProtocolNetworkID {
		log.Info().Msg("Downloader: InitEphemeryNetwork starting")
		if Workload.ClientName == "lighthouse" {
			ephemery_reset.ExtractAndDecEphemeralTestnetConfig(Workload.DataDir, beacon_cookbooks.LighthouseEphemeral)
		}
		if Workload.ClientName == "geth" {
			ephemery_reset.ExtractAndDecEphemeralTestnetConfig(Workload.DataDir, beacon_cookbooks.GethEphemeral)
		}
		log.Info().Msg("Downloader: InitEphemeryNetwork done")
	}
}
