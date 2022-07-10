package beacon_models

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/internal/beacon-api/api_types"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

func readBeaconValidatorsJSON() (api_types.ValidatorsStateBeacon, error) {
	var vsb api_types.ValidatorsStateBeacon
	jsonFile, err := os.Open("/Users/alex/Desktop/Zeus/olympus/configs/validators.json")
	if err != nil {
		return vsb, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &vsb)
	return vsb, err
}

func (vs *Validators) InsertValidatorsFromBeaconAPI(ctx context.Context) error {
	vsb, err := readBeaconValidatorsJSON()
	if err != nil {
		return err
	}
	apiValidatorStates := ToBeaconModelFormat(vsb)
	tmp := Validators{}
	for _, valState := range apiValidatorStates.Validators {
		if len(valState.Pubkey) >= 0 {
			val := Validator{Index: valState.Index,
				Pubkey:                     valState.Pubkey,
				Balance:                    valState.Balance,
				EffectiveBalance:           valState.EffectiveBalance,
				ActivationEligibilityEpoch: valState.ActivationEligibilityEpoch,
				ActivationEpoch:            valState.ActivationEpoch,
				ExitEpoch:                  valState.ExitEpoch,
				WithdrawableEpoch:          valState.WithdrawableEpoch,
				Slashed:                    valState.Slashed,
				WithdrawalCredentials:      valState.WithdrawalCredentials,
			}
			tmp.Validators = append(tmp.Validators, val)
		}

		if len(tmp.Validators) > 0 && len(tmp.Validators)%100 == 0 {
			err = InsertValidatorsFromBeaconAPI(ctx, tmp)
			tmp = Validators{}
		}
		time.Sleep(time.Second)
	}
	err = InsertValidatorsFromBeaconAPI(ctx, tmp)
	return err
}

var insertValidatorsFromBeaconAPI = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed, withdrawal_credentials) VALUES `

func InsertValidatorsFromBeaconAPI(ctx context.Context, vs Validators) error {
	vs.RowSetting.RowsToInclude = "beacon_state"
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorsFromBeaconAPI, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}
