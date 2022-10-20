package create

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Deployment struct {
	workloads.Deployment
}

func newDeployment() Deployment {
	return Deployment{workloads.NewDeployment()}
}

const ModelName = "Deployment"

func (d *Deployment) insertDeploymentStatement(c Chart) string {
	sqlInsertStatement := fmt.Sprintf(
		`%s, cte_insert_cct AS (
				    INSERT INTO chart_subcomponent_child_class_types(chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_name)
					VALUES ((SELECT chart_subcomponent_parent_class_type_id FROM cte_insert_pc), '%s')
				    RETURNING chart_subcomponent_child_class_type_id
				) SELECT chart_subcomponent_parent_class_type_id FROM cte_insert_pc
	`, d.insertDeploymentParentClass(c.ChartPackageID), "deploymentSpec")
	return sqlInsertStatement
}

func (d *Deployment) InsertDeployment(ctx context.Context, q sql_query_templates.QueryParams, c Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	r, err := apps.Pg.Exec(ctx, d.insertDeploymentStatement(c))
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}
