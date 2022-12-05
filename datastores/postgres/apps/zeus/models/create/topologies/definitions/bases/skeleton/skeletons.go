package create_skeletons

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "SkeletonBase"

func insertSkeletonBaseQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertSkeletonBaseDefinition"
	q.RawQuery = `INSERT INTO topology_skeleton_base_components (org_id, topology_base_component_id, topology_class_type_id, topology_skeleton_base_name)
			      VALUES ($1, $2, $3, $4)
				  RETURNING topology_skeleton_base_id`
	return q
}
func InsertSkeletonBase(ctx context.Context, skeletonBase *skeletons.SkeletonBase) error {
	q := insertSkeletonBaseQ()
	log.Debug().Interface("InsertSkeletonBase:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, skeletonBase.OrgID, skeletonBase.TopologyBaseComponentID, class_types.SkeletonBaseClassTypeID, skeletonBase.TopologySkeletonBaseName).Scan(&skeletonBase.TopologySkeletonBaseID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
