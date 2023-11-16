package eth_validators_service_requests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	olympus_hydra_validators_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/validators"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aegis_aws_secretmanager "github.com/zeus-fyi/zeus/pkg/aegis/aws/secretmanager"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	pods_client "github.com/zeus-fyi/zeus/zeus/z_client/workloads/pods"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
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
		log.Err(err)
		return err
	}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) error {
	var cloudCtxNs zeus_common_types.CloudCtxNs
	switch params.ProtocolNetworkID {
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		cloudCtxNs = EphemeryStakingCloudCtxNs
	case hestia_req_types.EthereumGoerliProtocolNetworkID:
		cloudCtxNs = GoerliStakingCloudCtxNs
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
		log.Err(err)
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
	case hestia_req_types.EthereumGoerliProtocolNetworkID:
		cloudCtxNs = GoerliStakingCloudCtxNs
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		cloudCtxNs = MainnetStakingCloudCtxNs
	default:
		return errors.New("unsupported protocol network id")
	}

	// TODO when capacity is reached, we need to delete the oldest pod
	//if vcNum > 0 {
	//	podName = fmt.Sprintf("%s-%d-0", olympus_hydra_validators_cookbooks.HydraValidatorsClientName, vcNum)
	//}
	vcNum := 0
	podName := fmt.Sprintf("%s-%d", olympus_hydra_validators_cookbooks.HydraValidatorsClientName, vcNum)
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			CloudCtxNs: cloudCtxNs,
		},
		Action:  zeus_pods_reqs.DeleteAllPods,
		PodName: podName,
	}
	pc := pods_client.NewPodsClientFromZeusClient(Zeus)
	_, err := pc.DeletePods(ctx, par)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

type Resty struct {
	*resty.Client
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) VerifyValidatorKeyOwnershipAndSigning(ctx context.Context, params ValidatorServiceGroupWorkflowRequest) (hestia_req_types.ValidatorServiceOrgGroupSlice, error) {
	totalVerifiedKeys := hestia_req_types.ValidatorServiceOrgGroupSlice{}
	feeAddrToPubkeyMap := make(map[string]string)

	totalKeys := len(params.ValidatorServiceOrgGroupSlice)
	var keyGroup hestia_req_types.ValidatorServiceOrgGroupSlice
	for _, vs := range params.ValidatorServiceOrgGroupSlice {
		feeAddrToPubkeyMap[vs.Pubkey] = vs.FeeRecipient
		keyGroup = append(keyGroup, vs)
		if len(keyGroup) >= 100 {
			newVerifiedKeys, retryKeys, verr := GetVerifiedKeys(ctx, feeAddrToPubkeyMap, params, keyGroup)
			if verr == nil {
				totalVerifiedKeys = append(totalVerifiedKeys, newVerifiedKeys...)
			} else {
				return nil, verr
			}
			time.Sleep(1 * time.Second)
			keyGroup = retryKeys
		}
	}
	if len(keyGroup) > 0 {
		newVerifiedKeys, _, verr := GetVerifiedKeys(ctx, feeAddrToPubkeyMap, params, keyGroup)
		if verr == nil {
			totalVerifiedKeys = append(totalVerifiedKeys, newVerifiedKeys...)
		} else {
			return nil, verr
		}
	}
	if len(totalVerifiedKeys) != totalKeys {
		return nil, errors.New("not all keys were verified")
	}
	return totalVerifiedKeys, nil
}

func GetVerifiedKeys(ctx context.Context, feeAddrToPubkeyMap map[string]string, params ValidatorServiceGroupWorkflowRequest, keyGroup hestia_req_types.ValidatorServiceOrgGroupSlice) ([]hestia_req_types.ValidatorServiceOrgGroup, []hestia_req_types.ValidatorServiceOrgGroup, error) {
	r := Resty{}
	r.Client = resty.New()
	req := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}

	for _, vs := range keyGroup {
		hexMessage, herr := aegis_inmemdbs.RandomHex(10)
		if herr != nil {
			panic(herr)
		}
		req.Map[vs.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{
			Message: hexMessage,
		}
	}
	sn := artemis_validator_signature_service_routing.FormatSecretNameAWS(params.ServiceRequestWrapper.GroupName, params.OrgID, params.ServiceRequestWrapper.ProtocolNetworkID)
	si := aegis_aws_secretmanager.SecretInfo{
		Region: awsSecretsRegion,
		Name:   sn,
	}
	sv, err := artemis_hydra_orchestrations_aws_auth.GetServiceRoutesAuths(ctx, si)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service routes auths")
		return nil, nil, err
	}
	signReqs := bls_serverless_signing.SignatureRequests{
		SecretName:        sv.ServiceAuth.SecretName,
		SignatureRequests: req,
	}
	respMsgMap := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)
	signedEventResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{
		Map: respMsgMap,
	}
	auth := aegis_aws_auth.AuthAWS{
		Region:    awsSecretsRegion,
		AccessKey: sv.ServiceAuth.AccessKey,
		SecretKey: sv.ServiceAuth.SecretKey,
	}
	reqAuth, err := auth.CreateV4AuthPOSTReq(ctx, "lambda", sv.ServiceAuth.ServiceURL, signReqs)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service routes auths for lambda iam auth")
		return nil, nil, err
	}
	r.SetBaseURL(sv.ServiceAuth.ServiceURL)
	// the first request make timeout, since it may have a cold start latency
	r.SetTimeout(12 * time.Second)
	r.SetRetryCount(5)
	r.SetRetryWaitTime(500 * time.Millisecond)
	resp, err := r.R().
		SetHeaderMultiValues(reqAuth.Header).
		SetResult(&signedEventResponses).
		SetBody(signReqs).Post("/")

	if err != nil {
		log.Err(err).Msg("failed to post to validator signature service url")
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Warn().Interface("respCode", resp.StatusCode()).Msg("resp code not 200, doing a cooldown")
		time.Sleep(10 * time.Second)
		return nil, keyGroup, nil
	}
	verifiedKeys, err := signedEventResponses.VerifySignatures(ctx, req, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to verify signatures")
		return nil, nil, err
	}
	vkAndFeeAddrSlice := make([]hestia_req_types.ValidatorServiceOrgGroup, len(verifiedKeys))
	for i, vk := range verifiedKeys {
		vkAndFeeAddrSlice[i] = hestia_req_types.ValidatorServiceOrgGroup{
			Pubkey:       vk,
			FeeRecipient: feeAddrToPubkeyMap[vk],
		}
	}
	return vkAndFeeAddrSlice, hestia_req_types.ValidatorServiceOrgGroupSlice{}, nil
}
