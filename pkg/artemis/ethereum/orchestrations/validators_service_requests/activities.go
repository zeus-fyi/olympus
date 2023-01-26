package eth_validators_service_requests

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
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
	return []interface{}{a.VerifyValidatorKeyOwnershipAndSigning, a.AssignValidatorsToCloudCtxNs, a.SendValidatorsToCloudCtxNs}
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

//func (z *ZeusClient) ReadDeployStatusUpdates(ctx context.Context, tar zeus_req_types.TopologyRequest) (zeus_resp_types.TopologyDeployStatuses, error) {
//	z.PrintReqJson(tar)
//	respJson := zeus_resp_types.TopologyDeployStatuses{}
//	resp, err := z.R().
//		SetResult(&respJson.Slice).
//		SetBody(tar).
//		Post(zeus_endpoints.DeployStatusV1Path)
//
//	if err != nil || resp.StatusCode() != http.StatusOK {
//		log.Ctx(ctx).Err(err).Msg("ZeusClient: ReadDeployStatusUpdates")
//		if err == nil {
//			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
//		}
//		return respJson, err
//	}
//	z.PrintRespJson(resp.Body())
//	return respJson, err
//}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params ArtemisEthereumValidatorsServiceRequestPayload) error {
	//err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, params.ProtocolNetworkID, params.CloudCtxNs.CloudProvider)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (a *ArtemisEthereumValidatorsServiceRequestActivities) SendValidatorsToCloudCtxNs(ctx context.Context, params ArtemisEthereumValidatorsServiceRequestPayload) error {
	// query all validators that should be in this cluster, then patch the validator yaml file

	return nil
}
