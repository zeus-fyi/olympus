package update_replicas

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (r *ReplicaUpdate) UpdateReplicaCountStatefulSet(ctx context.Context, replicaCount string) error {
	var q sql_query_templates.QueryParams
	q.QueryName = "UpdateReplicaCountStatefulSet"
	q.CTEQuery.Params = append(q.CTEQuery.Params, r.TopologyID, r.OrgID, r.UserID)
	log.Debug().Interface("UpdateReplicaCountStatefulSet:", q.LogHeader(Sn))
	updated := false
	q = UpdateReplicaCountSQL(q, replicaCount)
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, q.CTEQuery.Params...).Scan(&updated)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
