package create_deployments

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Deployment struct {
	deployments.Deployment
}

const ModelName = "Deployment"

func (d *Deployment) InsertDeployment(ctx context.Context, q sql_query_templates.QueryParams, c *charts.Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	q.CTEQuery = d.InsertDeploymentCte(c)
	q.RawQuery = q.CTEQuery.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Deployment: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}
