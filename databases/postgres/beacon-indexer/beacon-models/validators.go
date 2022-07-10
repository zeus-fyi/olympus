package beacon_models

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type Validator struct {
	Index                      int64
	Pubkey                     string
	Balance                    *int64
	EffectiveBalance           *int64
	ActivationEligibilityEpoch *uint64
	ActivationEpoch            *uint64
	ExitEpoch                  *uint64
	WithdrawableEpoch          *uint64
	Slashed                    bool
	Status                     *string
	WithdrawalCredentials      *string

	SubStatus  *string
	Network    string
	RowSetting ValidatorRowValuesForQuery
}

type Validators struct {
	Validators []Validator

	RowSetting ValidatorRowValuesForQuery
}

type ValidatorRowValuesForQuery struct {
	RowsToInclude string
}

func (v *Validator) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index, v.Pubkey}
	switch v.RowSetting.RowsToInclude {
	case "all":
		if v.Balance != nil && v.EffectiveBalance != nil && v.ActivationEligibilityEpoch != nil && v.ActivationEpoch != nil {
			pgValues = postgres.RowValues{v.Index, v.Pubkey, *v.Balance, *v.EffectiveBalance, *v.ActivationEligibilityEpoch, *v.ActivationEpoch}
		}
	default:
		pgValues = postgres.RowValues{v.Index, v.Pubkey}
	}
	return pgValues
}

func (vs *Validators) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, v := range vs.Validators {
		if len(vs.RowSetting.RowsToInclude) > 0 {
			v.RowSetting.RowsToInclude = vs.RowSetting.RowsToInclude
		}
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

func (v *Validator) getIndexValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index}
	return pgValues
}

func (v *Validator) getBalanceValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Balance}
	return pgValues
}
func (v *Validator) getEffectiveBalanceValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.EffectiveBalance}
	return pgValues
}

func (v *Validator) getActivationEligibilityEpochValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.ActivationEligibilityEpoch}
	return pgValues
}

func (v *Validator) getActivationEpochValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.ActivationEpoch}
	return pgValues
}

// GetManyRowValuesFlattened is just getting Index values for now
func (vs *Validators) GetManyRowValuesFlattened() postgres.RowValues {
	var pgRows postgres.RowValues
	for _, v := range vs.Validators {
		pgRows = append(pgRows, v.getIndexValues()...)
	}

	return pgRows
}

var insertValidatorsOnlyIndexPubkey = `INSERT INTO validators (index, pubkey) VALUES `

func (vs *Validators) InsertValidatorsOnlyIndexPubkey(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorsOnlyIndexPubkey, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

var insertValidatorPendingQueue = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch) VALUES `

func (vs *Validators) InsertValidatorsPendingQueue(ctx context.Context) error {
	vs.RowSetting.RowsToInclude = "all"
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorPendingQueue, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

func (vs *Validators) SelectValidatorsPendingQueue(ctx context.Context) (*Validators, error) {
	limit := 10000
	validators := strings.AnyArraySliceStrBuilderSQL(vs.GetManyRowValuesFlattened())
	query := fmt.Sprintf(`SELECT index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, status, substatus FROM validators WHERE index = %s LIMIT %d`, validators, limit)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEligibilityEpoch, &v.ActivationEpoch, &v.Status, &v.SubStatus)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, v)
	}
	return &selectedValidators, nil
}

func (vs *Validators) SelectValidatorsOnlyIndexPubkey(ctx context.Context) (*Validators, error) {
	limit := 10000
	validators := strings.AnyArraySliceStrBuilderSQL(vs.GetManyRowValuesFlattened())
	query := fmt.Sprintf(`SELECT index, pubkey FROM validators WHERE index = %s LIMIT %d`, validators, limit)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidators Validators
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index, &validator.Pubkey)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, validator)
	}
	return &selectedValidators, nil
}

func (vs *Validators) SelectValidators(ctx context.Context) (*Validators, error) {
	limit := 10000
	validators := strings.AnyArraySliceStrBuilderSQL(vs.GetManyRowValuesFlattened())
	query := fmt.Sprintf(`SELECT index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed, status, withdrawal_credentials, substatus FROM validators WHERE index = %s LIMIT %d`, validators, limit)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEpoch, &v.ActivationEpoch, &v.ExitEpoch, &v.WithdrawableEpoch, &v.Slashed, &v.WithdrawalCredentials, &v.SubStatus)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, v)
	}
	return &selectedValidators, nil
}

func (vs *Validators) UpdateValidatorBalancesAndActivationEligibility(ctx context.Context) (*Validators, error) {
	validators := strings.MultiArraySliceStrBuilderSQL(vs.GetManyRowValues())
	query := fmt.Sprintf(`
	WITH validator_update AS (
		SELECT * FROM UNNEST(%s) AS x(index, balance, effective_balance, activation_eligibility_epoch, activation_epoch)
	) 
	UPDATE validators
	SET balance = validator_updatebalance, effective_balance = validator_update.effective_balance, activation_eligibility_epoch = validator_update.activation_eligibility_epoch, activation_epoch = validator_update.activation_epoch
	JOIN validators ON validator_update.index = validators.index`, validators)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEpoch, &v.ActivationEpoch, &v.ExitEpoch, &v.WithdrawableEpoch, &v.Slashed, &v.WithdrawalCredentials, &v.SubStatus)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, v)
	}
	return &selectedValidators, nil
}
