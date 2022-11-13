package create_kns

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "TopologyKubeCtxNs"

func (k *Kns) InsertKns(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	r, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("TopologyKubeCtxNs: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func getInsertKnsQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertKns"

	q.RawQuery = `
	INSERT INTO topologies_kns(topology_id, cloud_provider, region, context, namespace, env)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT ON CONSTRAINT kns_pk DO NOTHING
	RETURNING true
	`
	return q
}

func InsertKns(ctx context.Context, kns *kns.TopologyKubeCtxNs) error {
	q := getInsertKnsQuery()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	success := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, kns.TopologyID, kns.CloudProvider, kns.Region, kns.Context, kns.Namespace, kns.Env).Scan(&success)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
