package create_skeletons

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func insertSkeletonBaseInfraQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertSkeletonBaseInfraDefinition"
	q.RawQuery = `INSERT INTO topology_skeleton_base_components (org_id, topology_base_component_id, topology_class_type_id, topology_skeleton_base_name)
			      VALUES ($1, $2, $3, $4)
				  RETURNING topology_skeleton_base_version_id`
	return q
}
