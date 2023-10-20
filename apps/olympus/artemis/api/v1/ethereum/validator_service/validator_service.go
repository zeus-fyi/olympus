package artemis_ethereum_validator_service

import (
	"context"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type DepositEthereumValidatorsService struct {
	Network               string                                        `json:"network"`
	ValidatorDepositSlice []signing_automation_ethereum.DepositDataJSON `json:"validatorDepositSlice"`
}

func CreateEthereumValidatorsHandler(c echo.Context) error {
	request := new(DepositEthereumValidatorsService)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DepositValidators(c)
}

func (v *DepositEthereumValidatorsService) DepositValidators(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	switch strings.ToLower(v.Network) {
	case "mainnet":
		return c.JSON(http.StatusNotImplemented, nil)
	case "goerli":
		token := c.Get("token").(string)
		key, rr := auth.VerifyBearerTokenService(ctx, token, create_org_users.EthereumGoerliService)
		if rr != nil || key.PublicKeyVerified == false {
			log.Err(rr).Interface("orgUser", ou).Msg("DepositValidators: EthereumGoerliService unauthorized")
			return c.JSON(http.StatusUnauthorized, nil)
		}
		resp := make([]ValidatorDepositResponse, len(v.ValidatorDepositSlice))
		w3client := signing_automation_ethereum.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumGoerli.NodeURL, artemis_network_cfgs.ArtemisEthereumGoerli.Account)
		h := make(map[string]string)
		h["Authorization"] = "Bearer " + token
		w3client.Headers = h

		txToBroadcast := make([]*types.Transaction, len(v.ValidatorDepositSlice))
		for i, d := range v.ValidatorDepositSlice {
			dp := signing_automation_ethereum.ExtendedDepositParams{
				ValidatorDepositParams: signing_automation_ethereum.ValidatorDepositParams{
					Pubkey:                d.Pubkey,
					WithdrawalCredentials: d.WithdrawalCredentials,
					Signature:             d.Signature,
					DepositDataRoot:       d.DepositDataRoot,
				},
				Amount:             d.Amount,
				DepositMessageRoot: d.DepositMessageRoot,
				ForkVersion:        d.ForkVersion,
				NetworkName:        hestia_req_types.Goerli,
			}
			signedTx, err := w3client.SignValidatorDepositTxToBroadcastFromJSON(ctx, dp)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("DepositValidators, goerli error")
				return c.JSON(http.StatusBadRequest, nil)
			}
			txToBroadcast[i] = signedTx
			rx, err := w3client.SubmitSignedTxAndReturnTxData(ctx, signedTx)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("DepositValidators, goerli error")
				return c.JSON(http.StatusBadRequest, nil)
			}
			resp[i] = ValidatorDepositResponse{
				Pubkey: d.Pubkey,
				Rx:     rx.Hash().String(),
			}
		}
		return c.JSON(http.StatusAccepted, resp)
	case "ephemery":
		resp := make([]ValidatorDepositResponse, len(v.ValidatorDepositSlice))
		w3client := signing_automation_ethereum.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumEphemeral.NodeURL, artemis_network_cfgs.ArtemisEthereumEphemeral.Account)
		txToBroadcast := make([]*types.Transaction, len(v.ValidatorDepositSlice))

		for i, d := range v.ValidatorDepositSlice {
			dp := signing_automation_ethereum.ExtendedDepositParams{
				ValidatorDepositParams: signing_automation_ethereum.ValidatorDepositParams{
					Pubkey:                d.Pubkey,
					WithdrawalCredentials: d.WithdrawalCredentials,
					Signature:             d.Signature,
					DepositDataRoot:       d.DepositDataRoot,
				},
				Amount:             d.Amount,
				DepositMessageRoot: d.DepositMessageRoot,
				ForkVersion:        d.ForkVersion,
				NetworkName:        hestia_req_types.Ephemery,
			}
			signedTx, err := w3client.SignValidatorDepositTxToBroadcastFromJSON(ctx, dp)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("DepositValidators, ephemery error")
				return c.JSON(http.StatusBadRequest, nil)
			}
			txToBroadcast[i] = signedTx
			rx, err := w3client.SubmitSignedTxAndReturnTxData(ctx, signedTx)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("DepositValidators, ephemery error")
				return c.JSON(http.StatusBadRequest, nil)
			}
			resp[i] = ValidatorDepositResponse{
				Pubkey: d.Pubkey,
				Rx:     rx.Hash().String(),
			}
		}
		return c.JSON(http.StatusAccepted, resp)
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
}

type ValidatorDepositResponse struct {
	Pubkey string `json:"pubkey"`
	Rx     string `json:"rx"`
}
