package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	age_encryption "github.com/zeus-fyi/zeus/pkg/crypto/age"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type GenerateValidatorsRequest struct {
	ValidatorDepositGenerationParams
}

func ValidatorsAgeEncryptedKeystoresGenerationRequestHandler(c echo.Context) error {
	request := new(GenerateValidatorsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GenerateValidatorsAgeEncryptedZipFile(c)
}

func (v *GenerateValidatorsRequest) GenerateValidatorsAgeEncryptedZipFile(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	vdg := signing_automation_ethereum.ValidatorDepositGenerationParams{
		Fp:                   filepaths.Path{},
		Mnemonic:             v.Mnemonic,
		Pw:                   v.HdWalletPw,
		ValidatorIndexOffset: v.HdOffset,
		NumValidators:        v.ValidatorCount,
		Network:              v.Network,
	}
	enc := age_encryption.NewAge(v.AgePrivKey, v.AgePubKey)
	inMemFs := memfs.NewMemFs()
	zip, err := vdg.GenerateAgeEncryptedValidatorKeysInMemZipFile(ctx, inMemFs, enc)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("GenerateValidatorsRequest, GenerateValidatorsAgeEncryptedZipFile error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, zip)
}

type ValidatorDepositGenerationParams struct {
	AgePubKey  string `json:"agePubKey,omitempty"`
	AgePrivKey string `json:"agePrivKey,omitempty"`

	Network        string `json:"network,omitempty"`
	Mnemonic       string `json:"mnemonic"`
	HdWalletPw     string `json:"hdWalletPw"`
	HdOffset       int    `json:"hdOffset"`
	ValidatorCount int    `json:"validatorCount"`
}
