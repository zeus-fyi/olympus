package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

var insertValidatorsOnlyIndexPubkey = `INSERT INTO validators (index, pubkey) VALUES `

func (vs *Validators) InsertValidatorsOnlyIndexPubkey(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorsOnlyIndexPubkey, vs.GetManyRowValues())
	r := postgres.Pg.QueryRow(ctx, query)
	err := r.Scan()
	if err != nil {
		return err
	}
	return err
}

var insertValidatorPendingQueue = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch) VALUES `

func (vs *Validators) InsertValidatorsPendingQueue(ctx context.Context) error {
	vs.RowSetting.RowsToInclude = "all"
	querySuffix := ` ON CONFLICT (index) DO NOTHING`
	query := strings.PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(insertValidatorPendingQueue, vs.GetManyRowValues(), querySuffix)
	r := postgres.Pg.QueryRow(ctx, query)
	err := r.Scan()
	if err != nil {
		return err
	}
	return err
}

func (vs *Validators) SelectValidatorsPendingQueue(ctx context.Context) (*Validators, error) {
	limit := 10000
	validators := strings.AnyArraySliceStrBuilderSQL(vs.GetManyRowValuesFlattened())
	query := fmt.Sprintf(`SELECT index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, status, substatus FROM validators WHERE index = %s LIMIT %d`, validators, limit)

	rows, err := postgres.Pg.Query(ctx, query)
	log.Err(err).Msg("SelectValidatorsPendingQueue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEligibilityEpoch, &v.ActivationEpoch, &v.Status, &v.SubStatus)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectValidatorsPendingQueue")
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
	log.Err(err).Msg("SelectValidatorsOnlyIndexPubkey")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var selectedValidators Validators
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index, &validator.Pubkey)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectValidatorsOnlyIndexPubkey")
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
	log.Err(err).Msg("SelectValidators")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEpoch, &v.ActivationEpoch, &v.ExitEpoch, &v.WithdrawableEpoch, &v.Slashed, &v.WithdrawalCredentials, &v.SubStatus)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectValidators")
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
	SET balance = validator_update.balance, effective_balance = validator_update.effective_balance, activation_eligibility_epoch = validator_update.activation_eligibility_epoch, activation_epoch = validator_update.activation_epoch
	JOIN validators ON validator_update.index = validators.index`, validators)

	rows, err := postgres.Pg.Query(ctx, query)
	log.Err(err).Msg("UpdateValidatorBalancesAndActivationEligibility")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var selectedValidators Validators
	for rows.Next() {
		var v Validator
		rowErr := rows.Scan(&v.Index, &v.Pubkey, &v.Balance, &v.EffectiveBalance, &v.ActivationEpoch, &v.ActivationEpoch, &v.ExitEpoch, &v.WithdrawableEpoch, &v.Slashed, &v.WithdrawalCredentials, &v.SubStatus)
		if rowErr != nil {
			log.Err(rowErr).Msg("UpdateValidatorBalancesAndActivationEligibility")
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, v)
	}
	return &selectedValidators, nil
}

func SelectCountValidatorEntries(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM validators`
	err := postgres.Pg.QueryRow(ctx, query).Scan(&count)
	log.Err(err).Msg("SelectCountValidatorEntries")
	if err != nil {
		return 0, err
	}
	return count, nil
}
