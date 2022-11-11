package containers

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateComputeResourcesCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	computeResourcesSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_compute_resources_%d", ts.UnixTimeStampNow()))
	computeResourcesSubCTE.TableName = "container_compute_resources"
	computeResourcesSubCTE.Columns = []string{"compute_resources_id", "compute_resources_cpu_request", "compute_resources_cpu_limit",
		"compute_resources_ram_request", "compute_resources_ram_limit", "compute_resources_ephemeral_storage_request", "compute_resources_ephemeral_storage_limit"}
	computeResourcesRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_compute_resources_relationships_%d", ts.UnixTimeStampNow()))
	computeResourcesRelationshipsSubCTE.TableName = "containers_compute_resources"
	computeResourcesRelationshipsSubCTE.Columns = []string{"compute_resources_id", "container_id"}
	return computeResourcesSubCTE, computeResourcesRelationshipsSubCTE
}
