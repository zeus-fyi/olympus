package olympus_hydra_validators_cookbooks

import (
	"context"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/consensus_client_apis/beacon_api"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
)

func GetValidatorsAndPrepareRemoteSignEmaulation(ctx context.Context, beaconURL, hydraAddres string) (ethereum_web3signer_actions.LighthouseWeb3SignerRequests, error) {
	state := "finalized"
	status := "active_ongoing"
	vs, err := beacon_api.GetValidatorsByState(ctx, beaconURL, state, status)
	if err != nil {
		return ethereum_web3signer_actions.LighthouseWeb3SignerRequests{}, err
	}
	lhw3 := make([]ethereum_web3signer_actions.LighthouseWeb3SignerRequest, len(vs.Data))
	req := ethereum_web3signer_actions.LighthouseWeb3SignerRequests{
		Enable:        true,
		Web3SignerURL: hydraAddres,
		FeeAddr:       "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
		Slice:         lhw3,
	}
	for i, v := range vs.Data {
		req.Slice[i] = ethereum_web3signer_actions.LighthouseWeb3SignerRequest{
			Enable:                true,
			SuggestedFeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
			VotingPublicKey:       v.Validator.Pubkey,
			Url:                   hydraAddres,
		}
	}
	return req, err
}
