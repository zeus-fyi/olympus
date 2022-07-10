package beacon_models

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

func SelectValidatorsToQueryBeacon(ctx context.Context, batchSize int) (*ValidatorBalancesEpoch, error) {
	query := fmt.Sprintf(`SELECT MAX(epoch) as max_epoch, validator_index
								 FROM validator_balances_at_epoch
								 WHERE epoch + 1 < (SELECT mainnet_finalized_epoch())
					             GROUP by validator_index
								 LIMIT %d`, batchSize)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidatorBalances ValidatorBalancesEpoch
	for rows.Next() {
		var vb ValidatorBalanceEpoch
		rowErr := rows.Scan(&vb.Epoch, &vb.Index)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidatorBalances.ValidatorBalance = append(selectedValidatorBalances.ValidatorBalance, vb)
	}
	return &selectedValidatorBalances, nil
}

func SelectValidatorIndexesInStrArrayForQueryURL(ctx context.Context, batchSize int) (string, error) {
	vbal, err := SelectValidatorsToQueryBeacon(ctx, batchSize)

	if err != nil {
		return "", err
	}
	var indexes []string
	indexes = make([]string, len(vbal.ValidatorBalance))
	for i, v := range vbal.ValidatorBalance {
		indexes[i] = strconv.FormatInt(v.Index, 10)

	}
	indexString := strings.UrlEncodeQueryParamList("", indexes...)
	return indexString, err
}
