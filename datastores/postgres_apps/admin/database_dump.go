package admin

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/cli_wrapper"
)

func DumpValidatorBalancesAtEpochTable(ctx context.Context, le, he int) (string, string, error) {
	subQuery := fmt.Sprintf("SELECT * FROM validator_balances_at_epoch WHERE epoch > %d AND epoch <= %d", le, he)
	return DumpTable(ctx, subQuery)
}

// TODO update the query statement
func DumpTable(ctx context.Context, subQuery string) (string, string, error) {
	execCmd := cli_wrapper.TaskCmd{}
	query := fmt.Sprintf(`psql -c "COPY (%s) TO STDOUT;" source_db | psql -c "COPY my_table FROM STDIN;" target_db`, subQuery)
	cmdStr := query
	execCmd.Command = cmdStr
	return execCmd.ExecuteCmd()
}
