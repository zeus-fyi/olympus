package skeletons

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "Skeleton"

type InsertSkeleton struct {
	skeletons.Skeleton
}

func insertSkeletonQ() sql_query_templates.QueryParams {

	q := sql_query_templates.QueryParams{}

	q.RawQuery = `INSERT INTO topology_skeleton_base_components (org_id, topology_base_component_id, topology_class_type_id, topology_skeleton_base_version_id, topology_skeleton_base_name)
			      VALUES ($1, $2)`

	return q
}
func (s *InsertSkeleton) InsertSkeleton(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertSkeleton:", q.LogHeader(Sn))
	r, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertTopologyClass: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
