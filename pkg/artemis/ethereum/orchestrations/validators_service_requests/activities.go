package eth_validators_service_requests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	olympus_hydra_validators_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/validators"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types/pods"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorsServiceRequestActivities struct {
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorsServiceRequestActivities {
	return ArtemisEthereumValidatorsServiceRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.ValidateKeysToServiceURL, d.AssignValidatorsToCloudCtxNs, d.RestartValidatorClient}
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) ValidateKeysToServiceURL(ctx context.Context, params hestia_autogen_bases.ValidatorServiceOrgGroup) error {
	client := Zeus
	client.BaseURL = params.ServiceURL
	// TODO needs to query for an auth token for external urls or run auth procedure, needs to validate message is signable
	// TODO POST random string payload, verify BLS key matches
	resp, err := client.R().Get("/health")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ArtemisEthereumValidatorsServiceRequestActivities: ValidateKeysToServiceURL health not-OK")
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return err
	}
	return nil
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) RestartValidatorClient(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	// this will pull the latest validators into the cluster
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			CloudCtxNs: params.CloudCtxNs,
		},
		Action:  zeus_pods_reqs.DeleteAllPods,
		PodName: fmt.Sprintf("%s-%d", olympus_hydra_validators_cookbooks.HydraValidatorsClientName, params.ValidatorClientNumber),
	}
	_, err := Zeus.DeletePods(ctx, par)
	if err != nil {
		return err
	}
	return nil
}
