package eth_validators_service_requests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	olympus_hydra_validators_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/validators"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types/pods"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorsServiceRequestActivities struct {
}

type ArtemisEthereumValidatorsServiceRequestPayload struct {
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorsServiceRequestActivities {
	return ArtemisEthereumValidatorsServiceRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{a.VerifyValidatorKeyOwnershipAndSigning, a.InsertVerifiedValidatorsWithFeeRecipient, a.AssignValidatorsToCloudCtxNs, a.RestartValidatorClient}
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) InsertVerifiedValidatorsWithFeeRecipient(
	ctx context.Context,
	params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol,
	pubkeys hestia_req_types.ValidatorServiceOrgGroupSlice) error {

	err := artemis_validator_service_groups_models.InsertVerifiedValidatorsToService(ctx, params, pubkeys)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	var cloudCtxNs zeus_common_types.CloudCtxNs
	switch params.ProtocolNetworkID {
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		cloudCtxNs = EphemeryStakingCloudCtxNs
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		cloudCtxNs = MainnetStakingCloudCtxNs
	default:
		return errors.New("unsupported protocol network id")
	}

	// TODO this needs to scope the pubkeys to the verified ones
	err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, params, cloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) RestartValidatorClient(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	// this will pull the latest validators into the cluster

	var cloudCtxNs zeus_common_types.CloudCtxNs
	switch params.ProtocolNetworkID {
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		cloudCtxNs = EphemeryStakingCloudCtxNs
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		cloudCtxNs = MainnetStakingCloudCtxNs
	default:
		return errors.New("unsupported protocol network id")
	}

	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			CloudCtxNs: cloudCtxNs,
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

type Resty struct {
	*resty.Client
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) VerifyValidatorKeyOwnershipAndSigning(ctx context.Context, params ArtemisEthereumValidatorsServiceRequestPayload) ([]string, error) {
	r := Resty{}
	r.Client = resty.New()
	req := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	for _, vs := range params.ValidatorServiceOrgGroupSlice {
		pubkey := vs.Pubkey
		req.Map[pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: rand.String(10)}
	}
	respJson := aegis_inmemdbs.EthereumBLSKeySignatureResponses{}
	_, err := r.R().
		SetResult(&respJson.Map).
		SetBody(req).
		Post(params.ServiceURL)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return nil, err
	}
	verifiedKeys, err := respJson.VerifySignatures(ctx, req)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return nil, err
	}
	return verifiedKeys, err
}
