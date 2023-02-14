package eth_validators_service_requests

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	olympus_hydra_validators_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/validators"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types/pods"
	"net/http"
)

const (
	awsSecretsRegion = "us-west-1"
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
	params artemis_validator_service_groups_models.OrgValidatorService,
	pubkeys hestia_req_types.ValidatorServiceOrgGroupSlice) error {

	err := artemis_validator_service_groups_models.InsertVerifiedValidatorsToService(ctx, params, pubkeys)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) error {
	var cloudCtxNs zeus_common_types.CloudCtxNs
	switch params.ProtocolNetworkID {
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		cloudCtxNs = EphemeryStakingCloudCtxNs
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		cloudCtxNs = MainnetStakingCloudCtxNs
	default:
		return errors.New("unsupported protocol network id")
	}
	// TODO when capacity is reached, we need to delete the oldest pod
	vsCloudCtx := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID:     params.ProtocolNetworkID,
		OrgID:                 params.OrgID,
		ValidatorClientNumber: 0,
	}
	err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, vsCloudCtx, cloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) RestartValidatorClient(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) error {
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

	// TODO when capacity is reached, we need to delete the oldest pod
	vcNum := 0
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			CloudCtxNs: cloudCtxNs,
		},
		Action:  zeus_pods_reqs.DeleteAllPods,
		PodName: fmt.Sprintf("%s-%d", olympus_hydra_validators_cookbooks.HydraValidatorsClientName, vcNum),
	}
	_, err := Zeus.DeletePods(ctx, par)
	if err != nil {
		log.Ctx(ctx).Error().Err(err)
		return err
	}
	return nil
}

type Resty struct {
	*resty.Client
}

// TODO add auth signature

func (a *ArtemisEthereumValidatorsServiceRequestActivities) VerifyValidatorKeyOwnershipAndSigning(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) (hestia_req_types.ValidatorServiceOrgGroupSlice, error) {
	r := Resty{}
	r.Client = resty.New()
	req := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	feeAddrToPubkeyMap := make(map[string]string)
	for _, vs := range params.ValidatorServiceOrgGroupSlice {
		pubkey := vs.Pubkey
		feeAddrToPubkeyMap[pubkey] = vs.FeeRecipient

		hexMessage, err := aegis_inmemdbs.RandomHex(10)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("unable to generate hex message")
		}
		req.Map[pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: hexMessage}
	}
	sn := artemis_validator_signature_service_routing.FormatSecretNameAWS(params.ServiceRequestWrapper.GroupName, params.OrgID, params.ServiceRequestWrapper.ProtocolNetworkID)
	si := aws_secrets.SecretInfo{
		Region: awsSecretsRegion,
		Name:   sn,
	}
	sv, err := artemis_hydra_orchestrations_aws_auth.GetServiceRoutesAuths(ctx, si)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get service routes auths")
		return nil, err
	}
	// TODO add auth signing on payload
	signReqs := bls_serverless_signing.SignatureRequests{
		SecretName:        sv.ServiceAuth.SecretName,
		SignatureRequests: req,
	}
	respMsgMap := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)
	signedEventResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{
		Map: respMsgMap,
	}
	resp, err := r.R().
		SetResult(&signedEventResponses).
		SetBody(signReqs).
		Post(sv.ServiceAuth.ServiceURL)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to post to validator signature service url")
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Warn().Msg("resp code not 200")
	}

	verifiedKeys, err := signedEventResponses.VerifySignatures(ctx, req, true)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to verify signatures")
		return nil, err
	}
	vkAndFeeAddrSlice := make([]hestia_req_types.ValidatorServiceOrgGroup, len(verifiedKeys))
	for i, vk := range verifiedKeys {
		vkAndFeeAddrSlice[i] = hestia_req_types.ValidatorServiceOrgGroup{
			Pubkey:       vk,
			FeeRecipient: feeAddrToPubkeyMap[vk],
		}
	}
	return vkAndFeeAddrSlice, nil
}
