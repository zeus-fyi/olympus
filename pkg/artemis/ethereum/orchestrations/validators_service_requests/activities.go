package eth_validators_service_requests

import (
	"context"
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

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorsServiceRequestActivities {
	return ArtemisEthereumValidatorsServiceRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{a.VerifyValidatorKeyOwnershipAndSigning, a.AssignValidatorsToCloudCtxNs}
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) RestartValidatorClient(ctx context.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
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

type ArtemisEthereumValidatorsServiceRequestPayload struct {
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice

	CloudCtxNs zeus_common_types.CloudCtxNs
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
