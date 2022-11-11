package containers

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateCmdArgsCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	cmdArgsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_command_args_%d", ts.UnixTimeStampNow()))
	cmdArgsSubCTE.TableName = "container_command_args"
	cmdArgsSubCTE.Columns = []string{"command_args_id", "command_values", "args_values"}
	cmdArgsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_command_args_relationships_%d", ts.UnixTimeStampNow()))
	cmdArgsRelationshipsSubCTE.TableName = "containers_command_args"
	cmdArgsRelationshipsSubCTE.Columns = []string{"command_args_id", "container_id"}
	return cmdArgsSubCTE, cmdArgsRelationshipsSubCTE
}
