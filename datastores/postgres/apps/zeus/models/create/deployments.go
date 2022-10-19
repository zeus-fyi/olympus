package create

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/code_templates/models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func insertDeploymentStatement(d workloads.Deployment) string {
	pc := d.ParentClassDefinition
	sqlInsertStatement := fmt.Sprintf(
		`WITH cte_insert AS (
					INSERT INTO chart_subcomponent_parent_class_types(chart_package_id, chart_component_resource_id, chart_subcomponent_parent_class_type_name)
					VALUES (%d, %d, %s
				)`, pc.ChartPackageID, pc.ChartComponentKindID, d.ParentClassDefinition.ChartSubcomponentParentClassTypeName)
	return sqlInsertStatement
}

func InsertDeployment(ctx context.Context, q sql_query_templates.QueryParams, d workloads.Deployment) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(models.Sn))
	r, err := apps.Pg.Exec(ctx, insertDeploymentStatement(d))
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(models.Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(models.Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(models.Sn))

}
