package beacon_models

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/cli_wrapper"
)

var pgConnStr = "postgresql://postgres:postgres@0.0.0.0:5432/postgres"

func DumpValidatorBalancesAtEpochTable(ctx context.Context, le, he int) (string, string, error) {
	subQuery := fmt.Sprintf("(SELECT * FROM validator_balances_at_epoch WHERE epoch > %d AND epoch <= %d);", le, he)
	csvName := fmt.Sprintf("validator_balances_at_epoch_%d_%d", le, he)
	return DumpTable(ctx, subQuery, csvName)
}

func DumpTable(ctx context.Context, subQuery, csvName string) (string, string, error) {
	execCmd := cli_wrapper.TaskCmd{Command: "psql"}
	cmdStr := fmt.Sprintf(`%s "%s" --csv postgres > %s.csv`, pgConnStr, subQuery, csvName)
	execCmd.Command = cmdStr
	return execCmd.ExecuteCmd()
}
