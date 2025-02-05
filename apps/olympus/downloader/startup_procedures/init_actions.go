package olympus_snapshot_init

import (
	"context"
	"path"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	beacon_cookbooks "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons"
	init_jwt "github.com/zeus-fyi/zeus/pkg/aegis/jwt"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/utils/ephemery_reset"
)

func InitWorkloadAction(ctx context.Context, w WorkloadInfo) {
	switch w.WorkloadType {
	case "validatorClient":
		log.Ctx(ctx).Info().Msg("starting validators sync")
		// TODO clientName is always lighthouse for validator clients for now, when you add others, add that conditional here
		err := w.DataDir.WipeDirIn()
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to wipe validator dir on startup")
			panic(err)
		}
		EphemeryReset(w)
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
		w.DataDir.FnOut = "validator_definitions.yml"
		w.DataDir.DirOut = path.Join(w.DataDir.DirIn, "/validators")
		err = w.DataDir.WriteToFileOutPath(ymlBytes)
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to write validators yaml")
			panic(err)
		}
		log.Info().Interface("validator_definitions.yml", lhW3Enable).Msg("validator_definitions.yml")
		log.Ctx(ctx).Info().Msg("validators sync complete")
	case "beaconExecClient", "beaconConsensusClient":
		EphemeryReset(w)
		log.Ctx(ctx).Info().Msgf("checking for upload snapshot job %s", w.WorkloadType)
		CheckForUploadBeaconJob(ctx, w)
		log.Ctx(ctx).Info().Msgf("checking for disk wipe job %s", w.WorkloadType)
		CheckForDiskWipeJobBeacon(ctx, w)
		log.Ctx(ctx).Info().Msg("starting chain sync")
		ChainDownload(ctx, w)
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

func EphemeryReset(w WorkloadInfo) {
	if w.ProtocolNetworkID == hestia_req_types.EthereumEphemeryProtocolNetworkID {
		log.Info().Msg("Downloader: InitEphemeryNetwork starting")
		if w.ClientName == "lighthouse" {
			ephemery_reset.ExtractAndDecEphemeralTestnetConfig(w.DataDir, beacon_cookbooks.LighthouseEphemeral)
		}
		if w.ClientName == "geth" {
			ephemery_reset.ExtractAndDecEphemeralTestnetConfig(w.DataDir, beacon_cookbooks.GethEphemeral)
		}
		log.Info().Msg("Downloader: InitEphemeryNetwork done")
	}
}
