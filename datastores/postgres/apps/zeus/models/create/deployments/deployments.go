package deployments

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Deployment struct {
	workloads.Deployment
}

func NewDeploymentConfigForDB(deployment workloads.Deployment) Deployment {
	return Deployment{deployment}
}

const ModelName = "Deployment"

func (d *Deployment) InsertDeployment(ctx context.Context, q sql_query_templates.QueryParams, c create.Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	r, err := apps.Pg.Exec(ctx, d.insertDeploymentCtes(c.ChartPackageID))
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}
