package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ethereum_automation_cookbook "github.com/zeus-fyi/zeus/cookbooks/ethereum/automation"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	age_encryption "github.com/zeus-fyi/zeus/pkg/crypto/age"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type GenerateValidatorsRequest struct {
	ValidatorDepositGenerationParams
}

func GenerateValidatorsHandler(c echo.Context) error {
	request := new(GenerateValidatorsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GenerateValidators(c)
}

func (v *GenerateValidatorsRequest) GenerateValidators(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	w3Client := signing_automation_ethereum.Web3SignerClient{
		Web3Actions: web3_actions.Web3Actions{
			NodeURL: "https://eth.ephemeral.zeus.fyi",
			Network: "ephemery",
		},
	}
	enc := age_encryption.NewAge(v.AgePrivKey, v.AgePubKey)
	vdg := signing_automation_ethereum.ValidatorDepositGenerationParams{
		Fp:                   filepaths.Path{},
		Mnemonic:             v.Mnemonic,
		Pw:                   v.HdWalletPw,
		ValidatorIndexOffset: v.HdOffset,
		NumValidators:        v.ValidatorCount,
		Network:              "ephemery",
	}
	// TODO needs to be a background job
	err := ethereum_automation_cookbook.GenerateValidatorDepositsAndCreateAgeEncryptedKeystores(ctx, w3Client, vdg, enc, v.HdWalletPw)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("GenerateValidatorsRequest, GenerateValidators error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

type ValidatorDepositGenerationParams struct {
	AgePubKey  string `json:"agePubKey"`
	AgePrivKey string `json:"agePrivKey"`

	Mnemonic       string `json:"mnemonic"`
	HdWalletPw     string `json:"hdWalletPw"`
	HdOffset       int    `json:"hdOffset"`
	ValidatorCount int    `json:"validatorCount"`
}
